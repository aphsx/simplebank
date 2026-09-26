// สร้างstore struct สืบทอดความสามารถของ Queries มาทั้งหมด
// และเพิ่มความสามารถในการเปิดTransaction จัดการ BEGIN,COMMIT, และ ROLLBACK

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
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx) //สร้างQueries ชุดใหม่ที่รันอยู่บนTransactionนี้
	err = fn(q)  // run callback funcที่ส่งเข้ามา
	if err != nil {
		//หากเกิด err ให้สั่งrollbackทันที
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	//หากคำส่ังทุกอย่างสำเร็จไม่มีerrให้สั่งบันทึกข้อมูล
	return tx.Commit()
}
