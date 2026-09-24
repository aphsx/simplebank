package main

// import (
// 	"errors"
// 	"fmt"
// )
// var ErrInsufficientBalance = errors.New("Balance is not enought")

// func Withdraw(balance int64, amount int64) (int64, error) {
// 	if amount <= 0 {
// 		return balance, errors.New("Withdraw should more than 0")
// 	}
// 	if amount > balance {
// 		return balance, ErrInsufficientBalance
// 	}
// 	return balance - amount, nil
// }

// func main() {
// 	newBalace, err := Withdraw(600, 100)
// 	if err != nil {
// 		fmt.Println("something went wrong")
// 		return
// 	}
// 	fmt.Println("Withdraw succes. Balance is", newBalace)
// }
