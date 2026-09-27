package db

//func TransferTx และเพ่ิมfunc addMoney กลไก Lock odering
import (
	"context"
	"database/sql"
	"fmt"
)

// คลุมทั้งfuncการqueryรายตัว และทำงานรวมกันเป็น transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create instancใหม่ของStore
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx จัดการTransaction กลางามี Auto Commit/Rollback
// execTx เป็น func คาบคุม Cycle ของTransaction
// db.BeginTx,tx.Rollback,tx.Commit ทุกครั้งที่ระบบจะทำอะไร
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	//1 ส่งให้dbเปิดTransaction (เริ่มACID)
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//2 create Query ชุดใหม่ทำงานใต้ Transaction (tx)
	q := New(tx)

	//3 run func ที่ส่งเข้ามา (withDraw,Transfers) โดยส่งQuery(q)ให้ไปใช้ต่อ
	err = fn(q)

	//4 check Func run is not Err
	if err != nil {
		//if err Rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}
	//if not Err Commit save
	return tx.Commit()
}
