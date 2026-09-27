package db

//func TransferTx และเพ่ิมfunc addMoney กลไก Lock odering
import (
	"context"
	"database/sql"
	"fmt"
)

// คลุมทั้งfuncการqueryรายตัว และทำงานรวมกันเป็น transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create instancใหม่ของStore
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx จัดการTransaction กลางามี Auto Commit/Rollback
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
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
