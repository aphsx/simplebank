package api

//api handler
//แกะJSON payload และ validatin ก่อนส่งไปบันทึกdb

import (
	"mybank/db"
	"net/http"
	db "simplebank/db"

	"gitgub.com/gin-gonic/gin"
)

//createAccountRequest กำหนด JSON ขาเข้า กับ (Binding tags) กฎการตรวจความถูกต้อง

type createAccountRequest struct {
	Owner    string `json:"owner binding:"required"`
	Currency string `json:"owner" binding: "required,oneof=USD EUR THB`
}

func (server *Server) createAccount(ctx *gin.Context) {
	//1 แกะJSON Payload ,and check Binding tags
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // ส่ง 400 Bad Request ถ้าผิดรูปแบบ
		return
	}

	//2
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	//3

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) //ส่ง500 if db มีปัญหา
		return
	}

	//4 ส่งข้อมูลAccount create succes back with JSON พร้อมhttp status 201 created
	ctx.JSON(http.StatusCreated, account)
}
