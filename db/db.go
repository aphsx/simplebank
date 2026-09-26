package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{db: tx}
}

// package db
// import (
// 	"database/sql"
// 	"time"
// 	_ "github.com/lib/pq" // postgres driver
// )

// func NewDB(dataSrouceName string) (*sql.DB,error){
// 	db,err:=sql.Open("postgres",dataSrouceName)
// 	if err != nil {
// 		return nil,err
// 	}
// 	db.SetMaxOpenConns(25)//limit connect
// 	db.SetMaxIdleConns(25) //limit connection ที่รอแสดนบายไว้ในpool
// 	db.SetConnMaxLifetime(5*time.Minute) //อายุุสูงสุดแต่ละconnection ก่อนจะ รีเฟรชใหม่

// 	//ทดสอบว่าเชื่อมฐานข้อมูลได้จริงหรือไม่
// 	if err:=db.Ping(); err != nil{
// 		return nil,err
// 	}
// 	return db,nil
// }
