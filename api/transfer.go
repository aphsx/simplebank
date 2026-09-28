package api

//ตรวจความถูกต้อง account from,account to and currency ก่อนเข้าTransaction
// Protect hackด้วยยอดเงินติดลบ

import (
	"database/sql"
	"fmt"
	"net/http"

	db "simplebank/db"

	"github.com/gin-gonic/gin"
)

// transferRequest กำหนดJSON for คำขอโอนเงิน แะเงื่อนไขความปลอดภัย
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	//1 แกะJSON payload and validate data
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//2 ตรวจสอบaccountต้นทางมีอยู่จริง และcurrencyตรง กับreq
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	//3check account to ปลายทาง มีอยู่จริง และcurrencyตรงกับ req
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	//4 เตรียม parameter ส่งต่อไป Transaction manager
	arg := db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	//5 run การโอนแบบ Atomic ่ผาน TransferTx
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//6 ส่งresult succes กลับไปหน้าบ้าน
	ctx.JSON(http.StatusOK, result)
}

//validAccount ตรวจก่อนเริ่มTransaction พร้อมคืนค่าAccount ไปใช้งานต่อ

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	//ค้นหา account from db
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err)) //ไม่เจอaccount 404
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	//checkควมสอดคล้องของcurrency ป้องกันการโอนข้ามสกุลเงิน

	if account.Currency != currency {
		err := fmt.Errorf("currency account[%d] is %s ไม่ตรงกับคำขอโอน %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
