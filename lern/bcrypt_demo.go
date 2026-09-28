package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	rawPassword := "secretPassword123"

	//1 hash same password 2 ครั้ง(Bcrypt สุ่มsalt ใหม่ทุกครั้ง ทำให้ผลลัพธ์ไม่ตรงกัน)
	hash1, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	hash2, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	fmt.Println("hash 1:", string(hash1))
	fmt.Println("hash 2:", string(hash2))
	fmt.Println("Hashes identical?:", string(hash1) == string(hash2))

	//2check password correct with CompareHashAndPadssword (ถ้าตรงกันerr=nil)
	errCorrect := bcrypt.CompareHashAndPassword(hash1, []byte("secretPassword123"))
	fmt.Println("Check correct password err == nil:", errCorrect == nil)

	//3 check passwordที่ผิด(if passwordไม่ตรงกัน จะส่งerr กลับมา)
	errWrong := bcrypt.CompareHashAndPassword(hash1, []byte("wrongPassword456"))
	fmt.Println("check wrong password error:", errWrong)
}
