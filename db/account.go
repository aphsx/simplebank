package db

//CRUD จัด table accounts ด้วย Parameterized Query
//แยกคำสั่งSQL ออกจากตัวแปรฐานข้อมูลป้องกัน SQL Injection ได้ 100%
import (
	"context"
)

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

const createAccount = ` -- name: CreateAccount :one
INSERT INTO accounts(
	owner,
	balance,
	currency,

)	
	VALUES(
	$1,$2,$3
)
RETURNING id,owner,balance,currency,create_at;
`

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreateAt,
	)
	return i, err
}

const getAccount = `-- namer: GetAccount :one
SELECT id,owner,
	balance,
	currency,
	create_at FROM accounts
WHERE id = $1 LIMIT 1;
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreateAt,
	)
	return i, err
}

type AddAccountBalanceParams struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"` //+ เงินเข้า,- เงินออก
}

const addAccountBalance = `-- name : AddAccountBalance :one
UPDATE accounts
SET balance = balance+$1
WHERE id = $2
RETURNING id,owner,balance,currency,create_at;
`

// AddAccountBalance ปรับยอดคงเหลือในระดับSQL เป็นAtomic (balance = balance + amount)
// รองรับทั้งการฝาก+,และการตัด -
func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	//1 UPDATE โดยให้ฐานข้อมูล+/-ยอดเงินในตัว(Atomic update)
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.ID)

	//2 แกะข้อมูลแถวที่updateแล้วใส่ตัวแปร Account
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreateAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	//1 DELTE ผ่านExecContextเมื่อไม่ต้องการข้อมูลแถวกลับมา
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

const listAccounts = `-- name: ListAccounts :many
SELECT id,owner,balance,currency,create_at FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2
`

//ListAccounts ดึงlist Accountsแบ่งหน้า (pagination)

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	//1 SELECT ค้นหาหลายแถวด้วย QueryContext
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//2 loopอ่านข้อมูลทีละแถวด้วย rows.Next()และ Scanค่าใส่Struct

	var items []Account

	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreateAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
