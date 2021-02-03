package database

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

// NovelStruct - sql structure for novel
type NovelStruct struct {
	ID        string    `db:"id"`
	Flags     string    `db:"flags"`
	Author    string    `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// GetNovels searches novel from given infomations
func GetNovels(db *sql.DB, id string, author string, all bool) []NovelStruct {
	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("*").From("novels").Where(
		builder.Or(
			builder.Equal("id", id),
			builder.Equal("author", author),
			builder.Equal("true", all)))

	sql, args := builder.Build()
	query, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}

	defer query.Close()
	var results []NovelStruct

	for query.Next() {
		var result NovelStruct
		err = query.Scan(&result.Content, &result.ID, &result.Author, &result.CreatedAt)

		if err != nil {
			panic(err)
		}

		results = append(results, result)
	}

	return results
}

// CheckNovelExists returns true when novel exists
func CheckNovelExists(db *sql.DB, id string, author string) bool {
	users := GetNovels(db, id, author, false)
	return len(users) > 0
}
