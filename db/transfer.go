//CRUD func save to table transfers
// save Entry 2 account พร้อมสร้าง Transfer ID

package db

import (
	"context"
)

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//

const createTransfer = `-- name :CreateTransfer :one
INSERT INTO transfers(
	from_account_id,
	to_account_id,
	amount
)	VALUES(
	$1,$2,$3
)

RETURNING id,from_account_id,to_account_id,amount,create_at;
`

//CreateTransfers บันทึกประวัติการโอนลง transfer trasfers

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Tranfer, error) {
	// ส่งคำสั่ง INSERT บันทึกการโอน
	rows := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)

	//2 แกะผลลัพธ์ Transfer
	var i Tranfer
	err := rows.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}
