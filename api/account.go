package api

//api handler
//แกะJSON payload และ validatin ก่อนส่งไปบันทึกdb

import (
	"database/sql"
	"net/http"
	db "simplebank/db"

	"simplebank/token"

	"github.com/gin-gonic/gin"
)

// createAccountRequest กำหนด JSON ขาเข้า กับ (Binding tags) กฎการตรวจความถูกต้อง
type createAccountRequest struct {
	// Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	//1 แกะJSON Payload ,and check Binding tags
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // ส่ง 400 Bad Request ถ้าผิดรูปแบบ
		return
	}

	//2
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	//ดึงข้อมูลเจ้าของtokenจากcontext
	arg := db.CreateAccountParams{
		// Owner:    req.Owner,
		Owner:    authPayload.Username,
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

// create getAccountID
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"requried,min1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			//ไม่เจอaccount 404 not found
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		//ฐานข้อมูลล่มตอบ 500 internal server error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account) //ส่งHTTP 200ok
}

// create listAccountsRequest
type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	//1 query parameters จาก URL (?page_id=1&page_size=5) ผ่าน ShouldBingQuery
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//2 calc Limit and Offset สำหรับส่งให้คำสั่งSQL (offset = (PageID -1 )*PageSize)
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	//3 ดึงlist accountจากdbผ่านstore
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//4 ส่งJSON  array accountกลับพร้อม HTTP 200 OK
	ctx.JSON(http.StatusOK, accounts)
}
