package endpoints

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pmh-only/ilggo.ga/src/database"
	"github.com/pmh-only/ilggo.ga/src/utils"
)

type commentCreationBody struct {
	Content string `json:"content"`
}

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
				"code":    334,
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

// CommuCommentListing returns comment
func CommuCommentListing(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    331,
				"success": false,
				"message": c.Param("id") + "는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    332,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		comments := database.GetComments(db, -1, id, false)
		c.JSON(200, gin.H{
			"code":    330,
			"success": true,
			"data":    comments,
		})
	}
}

// CommuCommentCreation auths user & create comment
func CommuCommentCreation(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		novelid, err := strconv.Atoi(c.Param("id"))

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    341,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    342,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		var body commentCreationBody
		err = c.ShouldBindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"code":    343,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		novels := database.GetNovels(db, novelid, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    344,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		if len(body.Content) < 1 {
			c.JSON(400, gin.H{
				"code":    345,
				"success": false,
				"message": "1자 미만의 댓글은 업로드 할 수 없습니다.",
			})
			return
		}

		id := 0
		comments := database.GetComments(db, -1, -1, true)

		if len(comments) > 0 {
			id = comments[0].ID + 1
		}

		database.CreateComment(db, id, novels[0].ID, user.ID, body.Content)
		c.JSON(201, gin.H{
			"code":    340,
			"success": true,
			"message": "댓글이 성공적으로 업로드 되었습니다.",
		})
	}
}

// CommuCommentUpdation checks permission and update novel infomations
func CommuCommentUpdation(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		cid, err := strconv.Atoi(c.Param("cid"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    351,
				"success": false,
				"message": c.Param("id") + " 혹은 " + c.Param("cid") + "는 올바른 자연수가 아닙니다.",
			})
			return
		}

		novels := database.GetNovels(db, id, "", false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    352,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if authHeader[0] != "Bearer" {
			c.JSON(401, gin.H{
				"code":    353,
				"success": false,
				"message": "인증 헤더가 잘못 입력되었습니다.",
			})
			return
		}

		user, err := utils.GetUsersFromJWT(db, authHeader[1], token)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    354,
				"success": false,
				"message": "토큰을 해석 할 수 없습니다.",
				"error":   err.Error(),
			})
			return
		}

		comments := database.GetComments(db, cid, -1, false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    352,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		if user.ID != comments[0].Author {
			c.JSON(403, gin.H{
				"code":    353,
				"success": false,
				"message": "본인이 쓴 댓글만 수정할 수 있습니다.",
				"error":   err.Error(),
			})
			return
		}

		var body commentCreationBody
		err = c.ShouldBindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"code":    354,
				"success": false,
				"message": "요청이 잘못되었습니다.",
			})
			return
		}

		if len(body.Content) < 1 {
			c.JSON(400, gin.H{
				"code":    355,
				"success": false,
				"message": "1자 미만의 댓글은 업로드 할 수 없습니다.",
			})
			return
		}

		database.UpdateComment(db, id, body.Content)

		c.JSON(200, gin.H{
			"code":    350,
			"success": true,
			"message": "댓글이 성공적으로 수정 되었습니다.",
		})
	}
}

// CommuCommentDeletion checks permissions and mark deleted
func CommuCommentDeletion(db *sql.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		cid, err := strconv.Atoi(c.Param("cid"))

		if err != nil || id < 0 {
			c.JSON(400, gin.H{
				"code":    241,
				"success": false,
				"message": c.Param("id") + " 혹은 " + c.Param("cid") + "는 올바른 자연수가 아닙니다.",
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

		comments := database.GetComments(db, cid, -1, false)
		if len(novels) < 1 {
			c.JSON(400, gin.H{
				"code":    242,
				"success": false,
				"message": "문서를 찾을 수 없습니다.",
			})
			return
		}

		if user.ID != comments[0].Author && user.ID != novels[0].Author {
			c.JSON(403, gin.H{
				"code":    245,
				"success": false,
				"message": "본인이 쓴 댓글만 수정할 수 있습니다.",
				"error":   err.Error(),
			})
			return
		}

		database.DeleteComment(db, id)

		c.JSON(200, gin.H{
			"code":    240,
			"success": true,
			"message": "댓글이 성공적으로 삭제 되었습니다.",
		})
	}
}
