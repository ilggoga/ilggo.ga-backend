package database

import (
	"database/sql"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

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
