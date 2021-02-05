package database

import (
	"database/sql"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

// CommentStruct - sql structure for comments
type CommentStruct struct {
	ID        int       `db:"id"`
	Novel     int       `db:"novel"`
	Author    string    `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// AddLikes adds user name into novel data's like feld
func AddLikes(db *sql.DB, id int, user string) {
	novels := GetNovels(db, id, "", false)

	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("novels").Where(builder.Equal("id", id)).Set(
			builder.Assign("likes", novels[0].Likes+user+","),
		).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}

// RemoveLikes removes user name into novel data's like feld
func RemoveLikes(db *sql.DB, id int, user string) {
	novels := GetNovels(db, id, "", false)
	newLikes := ""

	for _, like := range strings.Split(novels[0].Likes, ",") {
		if like == user {
			continue
		}

		newLikes += like + ","
	}

	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("novels").Where(builder.Equal("id", id)).Set(
			builder.Assign("likes", newLikes),
		).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}

// GetComments searches comment from given infomations
func GetComments(db *sql.DB, id int, novel int, all bool) []CommentStruct {
	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("*").From("comments").Where(
		builder.Or(
			builder.Equal("id", id),
			builder.Equal("novel", novel),
			builder.Equal("true", all),
		)).Desc().OrderBy("created_at")

	sql, args := builder.Build()
	query, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}

	defer query.Close()
	var results []CommentStruct

	for query.Next() {
		var result CommentStruct
		err = query.Scan(&result.ID, &result.Novel, &result.Author, &result.Content, &result.CreatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, result)
	}

	return results
}

// CreateComment creates comment and returns nothing
func CreateComment(db *sql.DB, id int, novel int, author string, content string) {
	builder := sqlbuilder.NewInsertBuilder()
	sql, args :=
		builder.InsertInto("comments").Cols("id", "novel", "author", "content").Values(id, novel, author, content).Build()

	_, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}
}

// UpdateComment updates comment infomations
func UpdateComment(db *sql.DB, id int, content string) {
	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("comments").Where(builder.Equal("id", id)).Set(
			builder.Assign("content", content),
		).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}

// DeleteComment deletes comment
func DeleteComment(db *sql.DB, id int) {
	builder := sqlbuilder.NewDeleteBuilder()
	sql, args :=
		builder.DeleteFrom("comments").Where(builder.Equal("id", id)).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}
