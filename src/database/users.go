package database

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

// UserStruct - sql structure for user
type UserStruct struct {
	ID        string         `db:"id"`
	Passwd    string         `db:"passwd"`
	Display   sql.NullString `db:"display"`
	IsAdmin   bool           `db:"is_admin"`
	CreatedAt time.Time      `db:"created_at"`
}

// CheckUserExists returns true when user exists
func CheckUserExists(db *sql.DB, id string, display string) bool {
	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("id").From("users").Or(builder.Equal("id", id), builder.Equal("display", display))

	sql, args := builder.Build()
	query, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}

	return query.Next()
}

// CreateUser creates user and returns nothing
func CreateUser(db *sql.DB, id string, display sql.NullString, passwdHash string) {
	builder := sqlbuilder.NewInsertBuilder()
	sql, args :=
		builder.InsertInto("users").Cols("id", "passwd", "display").Values(id, passwdHash, display).Build()

	_, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}
}
