package endpoints

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
	"github.com/pmh-only/ilggo.ga/src/utils"
)

type novelCreationBody struct {
	Flags   []string `json:"flags"`
	Content string   `json:"content"`
}

// NovelFetching checks novel exist & returns novel infomation
func NovelFetching(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    211,
				"success": false,
				"message": c.Param("id") + "는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    212,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		c.JSON(200, gin.H{
			"code":    210,
			"success": true,
			"data":    novels[0],
		})
	}
}

// NovelListing returns full list of novels
func NovelListing(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		novels := database.GetNovels(db, -1, "", true)
		c.JSON(200, gin.H{
			"code":    220,
			"success": true,
			"data":    novels,
		})
	}
}

// NovelCreation auths user & create novel content
func NovelCreation(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    231,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    232,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		var body novelCreationBody
		err = c.ShouldBindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"code":    233,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		if len(body.Content) < 50 {
			c.JSON(400, gin.H{
				"code":    234,
				"success": false,
				"message": "50자 미만의 문서는 업로드 할 수 없습니다.",
			})
			return
		}

		flagsStr := ""
		for _, flag := range body.Flags {
			flagsStr += flag + ","
		}

		id := 0
		novels := database.GetNovels(db, -1, "", true)

		if len(novels) > 0 {
			id = novels[0].ID + 1
		}

		database.CreateNovel(db, id, user.ID, body.Content, flagsStr)

		c.JSON(201, gin.H{
			"code":    230,
			"success": true,
			"message": "문서가 성공적으로 업로드 되었습니다.",
		})
	}
}

// NovelUpdation checks permission and update novel infomations
func NovelUpdation(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    241,
				"success": false,
				"message": c.Param("id") + "는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    242,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    243,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    244,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		if user.ID != novels[0].Author {
			c.JSON(403, gin.H{
				"code":    245,
				"success": false,
				"message": "본인이 쓴 게시글만 수정할 수 있습니다.",
				"error":   err.Error(),
			})
			return
		}

		var body novelCreationBody
		err = c.ShouldBindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"code":    246,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		if len(body.Content) < 50 {
			c.JSON(400, gin.H{
				"code":    247,
				"success": false,
				"message": "50자 미만의 문서는 업로드 할 수 없습니다.",
			})
			return
		}

		flagsStr := ""
		for _, flag := range body.Flags {
			flagsStr += flag + ","
		}

		database.UpdateNovel(db, id, body.Content, flagsStr)

		c.JSON(200, gin.H{
			"code":    240,
			"success": true,
			"message": "문서가 성공적으로 수정 되었습니다.",
		})
	}
}

// NovelDeletion checks permissions and mark deleted
func NovelDeletion(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    251,
				"success": false,
				"message": c.Param("id") + "는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    252,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    253,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    254,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		if user.ID != novels[0].Author {
			c.JSON(403, gin.H{
				"code":    255,
				"success": false,
				"message": "본인이 쓴 게시글만 삭제할 수 있습니다.",
				"error":   err.Error(),
			})
			return
		}

		database.DeleteNovel(db, id)
		c.JSON(200, gin.H{
			"code":    250,
			"success": true,
			"message": "문서가 성공적으로 삭제 되었습니다.",
		})
	}
}
