package api

import (
	db "simplebank/db"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//server จัดการHTTP requests ทั้งหมดของระบบ

type Server struct {
	store  *db.Store   //ที่เก็บตัวเชื่อมdbและTransaction
	router *gin.Engine // จัการ Endpoint HTTPS Routing ของ gin
}

//NewServer กำหนดค่าเริ่มต้น สร้างRouter และ Register Endpoint APIทั้งหมด

func NewServer(store *db.Store) *Server {
	//1สร้างInstance ของserver แนบ storeเข้าไป

	server := &Server{store: store}
	router := gin.Default()

	//register custom validator tag "currency"
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", varlidCurrency)
	}

	//2Register Routing แต่ละEndpoint
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)

	//3แนบRouterกลับเข้าไปในserver
	server.router = router
	return server
}

// start เปิดรันserverบนAddress ที่กำหนด (":8080")
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorRespone จัดรูปแบบerrให้อยู่ในรูปJson key "error"
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
