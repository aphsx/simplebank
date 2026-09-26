// สร้างstore struct สืบทอดความสามารถของ Queries มาทั้งหมด
// และเพิ่มความสามารถในการเปิดTransaction จัดการ BEGIN,COMMIT, และ ROLLBACK

//ทำเพื่อแก้ปัญหา เวลาเขียนคำสั่งQueryหลายๆตัวเรียงกัน เราต้องการระบบเปิดปิด txมาครอบมันไว้
// เพื่อความปลอดภ้ยของข้อมูล

package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store ครอบคลุมทั้งfunc การ queryรายตัว และการทำงานร่วมกันเป็น Transaction

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx เป็นฟังก์ชันภายในที่ควบคุมCycleของtransaction อย่างปลอดภ้ย
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//1สั่งให้ฐานข้อมูลเปิดTransaction (เริ่มหระบวนการคมข้อมูลแบบACID)
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//2 สร้างตัวQueryชุดใหม่ q โดยให้ทำงานอยู่ภายใต้ Transaction (tx) ที่เพิ่งเปิดไปด้านบน
	q := New(tx)
	//รันfunc คำสั่งที่เราส่งเข้ามา(เช่น Withdraw tranfers) โดยส่งตัวQuery (q) ให้มันไปใช้รันต่อ
	err = fn(q)
	//4เช็คว่าfuncที่รันไปไม่มีerr
	if err != nil {
		//ถ้า มีerr ให้rollbackทันที เพื่อยกเลิกทุกอย่างที่ทำไปในtxทันที
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err //ส่งerrออกไปบอกตัวเรียกใช้งาน
	}

	//5ถ้าไม่มี errเลยแปลว่าผ่าน ให้สั่งcommit เพื่อบันทึกข้อมูลลงไปถารวร
	return tx.Commit()
}
