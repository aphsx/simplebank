package main

import (
	"database/sql"
	"log"

	"simplebank/api"
	db "simplebank/db"
	"simplebank/util"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {
	// conn, err := sql.Open(dbDriver, dbSource)

	//1load setting from file app.env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	//2 coneect db from config
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db :", err)
	}

	store := db.NewStore(conn)
	// server := api.NewServer(store)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("connot create server :", err)
	}

	// err = server.Start(serverAddress)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
