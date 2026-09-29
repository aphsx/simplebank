//func จัดการบันทึกข้อมูล ledger table entries
//save เงินเข้าออก

package db

import (
	"context"
)

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"` //-ออก,+เข้า
}

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries(
	account_id,
	amount
)	VALUES(
	$1,$2
)
RETURNING id,account_id,amount,created_at;
`

// CreateEntry เข้าออก
func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)

	//2 แกะผลลัพธ์ใส่ตัวแปร Entry
	var i Entry

	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err

}
