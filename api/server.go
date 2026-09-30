package api

import (
	"fmt"
	db "simplebank/db"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//server จัดการHTTP requests ทั้งหมดของระบบ

type Server struct {
	config     util.Config
	store      *db.Store //ที่เก็บตัวเชื่อมdbและTransaction
	tokenMaker token.Maker
	router     *gin.Engine // จัการ Endpoint HTTPS Routing ของ gin
}

//NewServer กำหนดค่าเริ่มต้น สร้างRouter และ Register Endpoint APIทั้งหมด

// func NewServer(store *db.Store) *Server {
// 	//1สร้างInstance ของserver แนบ storeเข้าไป

// server := &Server{store: store}
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	//register custom validator tag "currency"
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", varlidCurrency)
	}

	//2Register Routing แต่ละEndpoint
	// router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccount)
	// router.GET("/accounts", server.listAccounts)
	// router.POST("/transfers", server.createTransfer)
	// router.POST("/users", server.createUser)

	//1เส้นpublic can access with out token
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	//private route ต้องผ่านmiddleware ก่อน
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.listAccounts)
	authRoutes.GET("/transfers", server.createTransfer)

	//3แนบRouterกลับเข้าไปในserver
	server.router = router
	return server, nil
}

// start เปิดรันserverบนAddress ที่กำหนด (":8080")
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorRespone จัดรูปแบบerrให้อยู่ในรูปJson key "error"
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
