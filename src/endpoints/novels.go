package endpoints

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
)

// NovelFetching checks novel exist & returns novel infomation
func NovelFetching(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    211,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		c.JSON(200, gin.H{
			"code":    210,
			"success": true,
			"data":    novels,
		})
	}
}
