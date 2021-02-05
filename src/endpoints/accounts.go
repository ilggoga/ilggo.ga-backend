package endpoints

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
)

type accountCreationBody struct {
	ID      string `json:"id"`
	Passwd  string `json:"passwd"`
	Display string `json:"display"`
}

type accountUpdationBody struct {
	ID        string `json:"id"`
	NewPasswd string `json:"newpasswd"`
	OldPasswd string `json:"oldpasswd"`
	Display   string `json:"display"`
}

type accountLoginBody struct {
	ID     string `json:"id"`
	Passwd string `json:"passwd"`
}

// AccountFetching checks account exist & returns account infomation
func AccountFetching(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		users := database.GetUsers(db, id, sql.NullString{Valid: true, String: ""})
		if len(users) < 1 {
			c.JSON(404, gin.H{
				"code":    141,
				"success": false,
				"message": "유저를 찾을 수 없습니다.",
			})
			return
		}

		users[0].Passwd = "<unknown>"

		c.JSON(200, gin.H{
			"code":    140,
			"success": true,
			"data":    users[0],
		})
	}
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

		display := sql.NullString{}
		if len(body.Display) > 0 {
			display = sql.NullString{
				Valid:  true,
				String: body.Display,
			}
		}

		exists := database.CheckUserExists(db, body.ID, display)
		if exists {
			c.JSON(409, gin.H{
				"code":    115,
				"success": false,
				"message": "해당 ID/닉네임은 이미 사용중입니다.",
			})
			return
		}

		hashFn := sha512.New()
		hashFn.Write([]byte(body.Passwd))

		passwdHash := hex.EncodeToString(hashFn.Sum(nil))
		database.CreateUser(db, body.ID, display, passwdHash)

		c.JSON(201, gin.H{
			"code":    110,
			"success": true,
			"message": "유저가 성공적으로 추가되었습니다.",
		})
	}
}

// AccountUpdation returns account updation endpoint
func AccountUpdation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body accountUpdationBody

		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    121,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		if len(body.ID) < 8 || len(body.ID) > 12 {
			c.JSON(400, gin.H{
				"code":    122,
				"success": false,
				"message": "ID가 올바르지 않은 형식입니다.",
			})
			return
		}

		if len(body.Display) > 50 {
			c.JSON(400, gin.H{
				"code":    123,
				"success": false,
				"message": "닉네임은 최대 50자 까지 사용가능 합니다.",
			})
		}

		if len(body.NewPasswd) < 8 {
			c.JSON(400, gin.H{
				"code":    124,
				"success": false,
				"message": "신규 비밀번호는 최소 8자리 이상이여야 합니다.",
			})
			return
		}

		users := database.GetUsers(db, body.ID, sql.NullString{Valid: true, String: ""})
		if len(users) < 1 {
			c.JSON(404, gin.H{
				"code":    125,
				"success": false,
				"message": "유저를 찾을 수 없습니다.",
			})
			return
		}

		displayUsers := database.GetUsers(db, "", sql.NullString{Valid: true, String: body.Display})
		if len(displayUsers) > 0 {
			if displayUsers[0].ID != body.ID {
				c.JSON(409, gin.H{
					"code":    126,
					"success": false,
					"message": "이미 사용중인 닉네임입니다.",
				})
				return
			}
		}

		hashFn := sha512.New()
		hashFn.Write([]byte(body.OldPasswd))

		oldPasswdHash := hex.EncodeToString(hashFn.Sum(nil))
		if users[0].Passwd != oldPasswdHash {
			c.JSON(403, gin.H{
				"code":    127,
				"success": false,
				"message": "기존 비밀번호가 일치하지 않습니다.",
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

		hashFn = sha512.New()
		hashFn.Write([]byte(body.NewPasswd))

		newPasswdHash := hex.EncodeToString(hashFn.Sum(nil))
		database.UpdateUser(db, body.ID, display, newPasswdHash)

		c.JSON(400, gin.H{
			"code":    120,
			"success": true,
			"message": "유저정보가 성공적으로 수정되었습니다.",
		})
	}
}

// AccountLogin checks account infomations & generate jwt
func AccountLogin(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body accountLoginBody

		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{
				"code":    131,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		users := database.GetUsers(db, body.ID, sql.NullString{Valid: true, String: ""})
		if len(users) < 1 {
			c.JSON(400, gin.H{
				"code":    132,
				"success": false,
				"message": "유저를 찾을 수 없습니다.",
			})
			return
		}

		hashFn := sha512.New()
		hashFn.Write([]byte(body.Passwd))

		passwdHash := hex.EncodeToString(hashFn.Sum(nil))
		if users[0].Passwd != passwdHash {
			c.JSON(403, gin.H{
				"code":    133,
				"success": false,
				"message": "비밀번호가 일치하지 않습니다.",
			})
			return
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  body.ID,
			"iat": time.Now().Unix(),
			"exp": time.Now().AddDate(0, 0, 1).Unix(),
		}).SignedString([]byte(token))

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"code":    130,
			"success": true,
			"message": "로그인이 완료되었습니다.",
			"data":    token,
		})
	}
}
