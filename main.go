package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/endpoints"
)

func main() {
	router := gin.Default()
	router.POST("/api/account", endpoints.AccountCreation())
	router.Static("/", "./dist")
	router.Run()
}
