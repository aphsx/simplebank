// นิยาม struct สำหรับพารามิเตอร์ขาเข้า และ ผลลัพธ์ ขาออกของธุรกรรมการโอนเงิิน
// มัดรวมข้อมูลที่เกี่ยวข้องเป้นเป้นก้อนเดียว เพื่อส่งผ่านและตรวจสอบข้อมูลได้ง่าย ปลอดภัย

//เขียนstructสำหรับ requestและ result

package db

import "context"

// tranferTxParams ข้อมูลที่ต้องใช้ในการสั่งโอน ระหว่าง2 acccount

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// tranferTxResult ผลลัพธ์ที่ได้หลังจากโอนเงินสำเร็จ เพื่อส่งกลับให้หน้าบ้านแสดงผล
type TransferTxResult struct {
	Transfer    Tranfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

//addMoney updata balane ของทั้ง2 account ตามลำดับที่ส่งเข้ามา

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	///1 update balance account1 (lock row accountID1 ก่อน)
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	//2 update balance account2 (lock row accountID2 ตามมา)
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

// TransferTx version Protect Deadlock
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//1 createประวัติการโอนลงtable transfers
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//2 create Entry เงินออกaccountต้นทาง (ledgerบันทึกยอดลบ)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		//3create Entry เงินเข้าAccountปลายทาง (ledgerบันทึกยอดเข้า+)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//4 Main ตรวจID ไหนน้อยกว่า ให้Update Accountนั้นก่อนเสมอ !
		// ไม่ว่าใครจะโอนหรือรับ ลำดับการขอlock ต้องวิ่งจากIDน้ยอไปIDมากเสมอ

		if arg.FromAccountID < arg.ToAccountID {
			//case1 from id < to id หักเงินบัญชีFromก่อน แล้วค่อยเพ่ิมเงินบัญชีto
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount,
			)

		} else {
			//case2 To id น้อยกว่า from idเพ่ิมเงินบััญชี Toก่อนแล้วค่อยหักFrom
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx, q, arg.ToAccountID, -arg.Amount, arg.FromAccountID, arg.Amount,
			)
		}
		return err

	})

	return result, err
}
