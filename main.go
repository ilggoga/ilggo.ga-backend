package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/endpoints"
)

func main() {
	db, err := sql.Open("mysql", "ilggoga@tcp(localhost:3306)/ilggoga?parseTime=true")
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.POST("/api/account", endpoints.AccountCreation(db))
	router.PUT("/api/account", endpoints.AccountUpdation(db))
	router.Static("/", "./dist")
	router.Run()
}
