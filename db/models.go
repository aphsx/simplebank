package db

import (
	"time"
)

type Account struct {
	ID int64 `json:"id"`
	Owner string `json:"owner"`
	Balance int64 `json:"balance"`
	Currency string `json:"currency"`
	CreateAt time.time `json:"create_at"`
}

type Entry struct {
	ID int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	Amount int64 `json:"amount"` // + เงินเข้า , - เงินออก
	CreateAt time.time `json:"create_at"`
}

type Tranfer struct {
	ID int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
	CreateAt time.time `json:"create_at"`
}