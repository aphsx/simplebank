package main

import (
	"errors"
	"fmt"
)

// ErrInsufficientBalance
var ErrInsufficientBalance = errors.New("Balance not enought")

func Withdraw(balance int64, amount int64) (int64, error) {
	if amount <= 0 {
		return balance, errors.New("Withdraw should more than 0")
	}

	if balance < amount {
		return balance, ErrInsufficientBalance
	}
	return balance - amount, nil
}

func main() {
	newBalance, err := Withdraw(500, 400)
	if err != nil {
		fmt.Println("Something went wrong")
		return
	}
	fmt.Println("Withdraw succes. Balance is: ", newBalance)
}
