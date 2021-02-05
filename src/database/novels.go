package database

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

// NovelStruct - sql structure for novel
type NovelStruct struct {
	ID        int       `db:"id"`
	Likes     string    `db:"likes"`
	Flags     string    `db:"flags"`
	Author    string    `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// GetNovels searches novel from given infomations
func GetNovels(db *sql.DB, id int, author string, all bool) []NovelStruct {
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
		err = query.Scan(&result.Content, &result.ID, &result.Author, &result.CreatedAt, &result.Flags, &result.Likes)

		if err != nil {
			panic(err)
		}

		results = append(results, result)
	}

	return results
}

// CheckNovelExists returns true when novel exists
func CheckNovelExists(db *sql.DB, id int, author string) bool {
	users := GetNovels(db, id, author, false)
	return len(users) > 0
}

// CreateNovel creates novel and returns nothing
func CreateNovel(db *sql.DB, id int, author string, content string, flags string) {
	builder := sqlbuilder.NewInsertBuilder()
	sql, args :=
		builder.InsertInto("novels").Cols("id", "author", "content", "flags").Values(id, author, content, flags).Build()

	_, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}
}

// UpdateNovel updates novel infomations
func UpdateNovel(db *sql.DB, id int, content string, flags string) {
	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("novels").Where(builder.Equal("id", id)).Set(
			builder.Assign("content", content),
			builder.Assign("flags", flags),
		).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}

// DeleteNovel deletes novel
func DeleteNovel(db *sql.DB, id int) {
	novels := GetNovels(db, id, "", false)
	flags := novels[0].Flags + "deleted,"

	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("novels").Where(builder.Equal("id", id)).Set(builder.Assign("flags", flags)).Build()

	_, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
}
