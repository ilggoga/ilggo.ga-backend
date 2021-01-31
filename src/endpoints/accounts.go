package endpoints

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
)

type accountCreationBody struct {
	ID      string `json:"id"`
	Passwd  string `json:"passwd"`
	Display string `json:"Display"`
}

// AccountCreation returns account endpoint
func AccountCreation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body accountCreationBody

		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    111,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		if len(body.ID) < 8 || len(body.ID) > 12 {
			c.JSON(400, gin.H{
				"code":    112,
				"success": false,
				"message": "ID는 8자에서 12자 사이여야 합니다.",
			})
			return
		}

		if len(body.Display) > 50 {
			c.JSON(400, gin.H{
				"code":    113,
				"success": false,
				"message": "닉네임은 최대 50자 까지 사용가능 합니다.",
			})
		}

		if len(body.Passwd) < 8 {
			c.JSON(400, gin.H{
				"code":    114,
				"success": false,
				"message": "비밀번호는 최소 8자리 이상이여야 합니다.",
			})
			return
		}

		hashFn := sha512.New()
		hashFn.Write([]byte(body.Passwd))

		passwdHash := hex.EncodeToString(hashFn.Sum(nil))
		exists := database.CheckUserExists(db, body.ID, body.Display)
		if exists {
			c.JSON(400, gin.H{
				"code":    115,
				"success": false,
				"message": "해당 ID/닉네임은 이미 사용중입니다.",
			})
			return
		}

		display := sql.NullString{}
		if len(body.Display) > 0 {
			display = sql.NullString{
				Valid:  true,
				String: body.Display,
			}
		}

		database.CreateUser(db, body.ID, display, passwdHash)

		c.JSON(201, gin.H{
			"code":    110,
			"success": true,
			"message": "유저가 성공적으로 추가되었습니다.",
		})
	}
}
