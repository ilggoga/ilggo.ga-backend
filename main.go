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

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "5 by 5",
		})
	})

	router.POST("/auth", endpoints.AccountLogin(db, token))
	router.GET("/accounts/:id", endpoints.AccountFetching(db))
	router.POST("/accounts", endpoints.AccountCreation(db))
	router.PUT("/accounts", endpoints.AccountUpdation(db))

	router.GET("/novels", endpoints.NovelListing(db))
	router.POST("/novels", endpoints.NovelCreation(db, token))
	router.GET("/novels/:id", endpoints.NovelFetching(db))
	router.PUT("/novels/:id", endpoints.NovelUpdation(db, token))
	router.DELETE("/novels/:id", endpoints.NovelDeletion(db, token))

	router.PUT("/novels/:id/like", endpoints.CommuAddLike(db, token))
	router.DELETE("/novels/:id/like", endpoints.CommuRemoveLike(db, token))
	router.GET("/novels/:id/comments", endpoints.CommuCommentListing(db))
	router.POST("/novels/:id/comments", endpoints.CommuCommentCreation(db, token))
	router.PUT("/novels/:id/comments/:cid", endpoints.CommuCommentUpdation(db, token))
	router.DELETE("/novels/:id/comments/:cid", endpoints.CommuCommentDeletion(db, token))

	router.Run()
}
