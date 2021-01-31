package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/endpoints"
	"github.com/pmh-only/ilggo.ga/src/utils"
)

func main() {
	token, err := utils.GenerateRandomString(30)
	db, err := sql.Open("mysql", "ilggoga@tcp(localhost:3306)/ilggoga?parseTime=true")
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/api/account", endpoints.AccountLogin(db, token))
	router.POST("/api/account", endpoints.AccountCreation(db))
	router.PUT("/api/account", endpoints.AccountUpdation(db))

	router.GET("/", func(c *gin.Context) { c.Redirect(301, "/bbs") })
	router.Static("/bbs", "./dist")
	router.Run()
}
