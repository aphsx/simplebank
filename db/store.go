// สร้าง store struct ที่สอบทอดความสามรถของ Queries มาทั้งหมด และเพ่ิมความสามารถในการเปิด  Transaction
// จัการBEGIN,COMMIT ,และROLLBACK
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store ครอบคลุมทั้งfuncการqueryรายตัว และทำงานรวมกันเป็นTransaction
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

//ecexTx funcภายในควลคุมรอบชีวติของ Transaction อย่างปลอดภัย

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx) // สร้างQueriesชุดใหม่ที่รันนอยู่บนTransaction นี้
	err = fn(q)
	if err != nil {
		//หากเกิด err ให้ย้อนกลับทั้งที
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	//หากคำส่ังทุกอย่างาสำเร็จ ไมมีerrใหสั่งยืนยนับันทึกข้อมูล
	return tx.Commit()
}
