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

// GetUsers searches user from given infomations
func GetUsers(db *sql.DB, id string, display sql.NullString) []UserStruct {
	builder := sqlbuilder.NewSelectBuilder()

	builder.Select("*").From("users").Where(builder.Or(builder.Equal("id", id), builder.Equal("display", display))).OrderBy("id")

	sql, args := builder.Build()
	query, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}

	defer query.Close()
	var results []UserStruct

	for query.Next() {
		var result UserStruct
		err = query.Scan(&result.ID, &result.Passwd, &result.Display, &result.IsAdmin, &result.CreatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, result)
	}

	return results
}

// CheckUserExists returns true when user exists
func CheckUserExists(db *sql.DB, id string, display sql.NullString) bool {
	users := GetUsers(db, id, display)
	return len(users) > 0
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

// UpdateUser updates user infomations
func UpdateUser(db *sql.DB, id string, display sql.NullString, passwdHash string) {
	builder := sqlbuilder.NewUpdateBuilder()
	sql, args :=
		builder.Update("users").Where(builder.Equal("id", id)).Set(
			builder.Assign("display", display),
			builder.Assign("passwd", passwdHash),
		).Build()

	_, err := db.Query(sql, args...)

	if err != nil {
		panic(err)
	}
}
