package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword แปลง pass Plain text เป้นBcrypt Hash + salt อัตนโมมัติ
func HashPassword(password string) (string, error) {
	//1 hashpassโดยใช้ cost มาตรฐาน (DefaultCost = 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถแฮชรหัสผ่านได้: %w", err)
	}

	//2คืนค่าstringของBcrypt Hashที่พร้อม save ลง คอลัมน์ hashed_password
	return string(hashedPassword), nil
}

// CheckPassword ที่ใส่เข้ามาตรงกบัBcrypt Hash ในdb
func CheckPassword(password string, hashedPassword string) error {
	//use CompareHashAndPassword แกะsaltจากHash และเทียบpassที่ส่งเข้ามา
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
