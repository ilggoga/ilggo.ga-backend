package endpoints

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
	"github.com/pmh-only/ilggo.ga/src/utils"
)

// CommuAddLike checks like conflict and add like
func CommuAddLike(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    311,
				"success": false,
				"message": c.Param("id") + "(은)는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    312,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    313,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    314,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		likedUsers := strings.Split(novels[0].Likes, ",")
		if utils.IsInStringSlice(likedUsers, user.ID) {
			c.JSON(409, gin.H{
				"code":    315,
				"success": false,
				"message": "이미 좋아요를 눌렀습니다.",
			})
			return
		}

		database.AddLikes(db, id, user.ID)
		c.JSON(201, gin.H{
			"code":    310,
			"success": true,
			"message": "좋아요가 반영되었습니다.",
		})
	}
}

// CommuRemoveLike checks like conflict and remove like
func CommuRemoveLike(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    321,
				"success": false,
				"message": c.Param("id") + "(은)는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    322,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    323,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    324,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		likedUsers := strings.Split(novels[0].Likes, ",")
		if !utils.IsInStringSlice(likedUsers, user.ID) {
			c.JSON(409, gin.H{
				"code":    325,
				"success": false,
				"message": "좋아요를 누른 기록이 없습니다.",
			})
			return
		}

		database.RemoveLikes(db, id, user.ID)
		c.JSON(201, gin.H{
			"code":    320,
			"success": true,
			"message": "좋아요 취소가 반영되었습니다.",
		})
	}
}
