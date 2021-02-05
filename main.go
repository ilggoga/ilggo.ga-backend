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

	router.POST("/api/auth", endpoints.AccountLogin(db, token))
	router.POST("/api/accounts", endpoints.AccountCreation(db))
	router.PUT("/api/accounts", endpoints.AccountUpdation(db))

	router.GET("/api/novels", endpoints.NovelListing(db))
	router.POST("/api/novels", endpoints.NovelCreation(db, token))
	router.GET("/api/novels/:id", endpoints.NovelFetching(db))
	router.PUT("/api/novels/:id", endpoints.NovelUpdation(db, token))
	router.DELETE("/api/novels/:id", endpoints.NovelDeletion(db, token))

	router.PUT("/api/novels/:id/like", endpoints.CommuAddLike(db, token))
	router.DELETE("/api/novels/:id/like", endpoints.CommuRemoveLike(db, token))

	router.GET("/", func(c *gin.Context) { c.Redirect(301, "/bbs") })
	router.Static("/bbs", "./dist")
	router.Run()
}
