// นิยาม struct สำหรับพารามิเตอร์ขาเข้า และ ผลลัพธ์ ขาออกของธุรกรรมการโอนเงิิน
// มัดรวมข้อมูลที่เกี่ยวข้องเป้นเป้นก้อนเดียว เพื่อส่งผ่านและตรวจสอบข้อมูลได้ง่าย ปลอดภัย

package db

// tranferTxParams ข้อมูลที่ต้องใช้ในการสั่งโอน ระหว่าง2 acccount

type TranferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// tranferTxResult ผลลัพธ์ที่ได้หลังจากโอนเงินสำเร็จ เพื่อส่งกลับให้หน้าบ้านแสดงผล

type TranferTxResult struct {
	Tranfer     Tranfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}
