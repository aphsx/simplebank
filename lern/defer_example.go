package main

import "fmt"

func ProcessTransaction() {
	fmt.Println("1. Start DB connect")
	defer fmt.Println("5. Close DB connect(working when func end)")
	fmt.Println("2. Check account")
	fmt.Println("3. Tranfer")
	fmt.Println("4. History save succes")
}

func main() {
	ProcessTransaction()
}
