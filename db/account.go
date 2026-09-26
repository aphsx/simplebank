package db

import (
	"context"
)


type CreateAccountParams struct{
	Owner string `json:"owner"`
	Balance int64 `json:"balance"`
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


func (q *Queries) CreateAccount(ctx context.Context,arg CreateAccountParams)(Account,error){
	row := q.db.QueryRowContext(ctx,createAccount,arg.Owner,arg.Balance,arg.Currency)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreateAt,
	)
	return i,err
}


const getAccount = `-- namer: GetAccount :one
SELECT id,owner,
	balance,
	currency,
	create_at FROM accounts
WHERE id = $1 LIMIT 1;
`
func (q *Queries) GetAccount(ctx context.Context,id int64)(Account,error){
	row := q.db.QueryRowContext(ctx,getAccount,id)
	var i accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreateAt,
	)
	return i,err
}