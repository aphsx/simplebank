//tranferTx รวม5 ขั้น ของการโอนเงินจริงเป็น1 Transaction

//ทุกขั้นตอนจะถูกครอบด้วย execTx ถ้าขั้นตอนไหนerrระบบจะRollback คืนค่าเดิม

package db

import (
	"context"
)

// TransferTx runการโอนเิงนจากบํญชีหนึ่งไปอีกบัญชีใน1 Transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// เรียกใช้ execTx เพื่อเปิดTransaction และครอบทุกขั้นตอนไว้ด้วยกัน
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//1สร้างแถวในตาราง transfers เพื่อบันทึกประวติการโอน

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//2สร้างentry เงินออกของบัญชีต้นทาง (ยอดเงินจะติดลบในLedger)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3สร้างentryเข้าบัญชีปลายทาง (ยอดเงินใน Ledgerเป็น+)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 4 หักเงินบัญชีต้นทาง
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		//5เพิ่มยอดบัญชีปลายทาง
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		//ทุกขั้นตอนsuccesไม่มีerrให้คืน nillเพื่อสั่งcommmit
		return nil
	})
	return result, err
}
