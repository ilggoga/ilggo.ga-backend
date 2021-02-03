package utils

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	"github.com/dgrijalva/jwt-go"
	"github.com/pmh-only/ilggo.ga/src/database"
)

// from https://gist.github.com/denisbrodbeck/635a644089868a51eccd6ae22b2eb800#file-main-go
func GenerateRandomString(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}

		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}

		n := num.Int64()
		// Make sure that the number/byte/letter is inside
		// the range of printable ASCII characters (excluding space and DEL)
		if n > 32 && n < 127 {
			result += string(rune(n))
		}
	}
}

// from https://stackoverflow.com/a/27272103
func IsInStringSlice(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

type jwtCtx struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// GetUsersFromJWT checks token is vaild and returns user info
func GetUsersFromJWT(db *sql.DB, token string, key string) (database.UserStruct, error) {
	proc, err := jwt.ParseWithClaims(token, &jwtCtx{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if claims, ok := proc.Claims.(*jwtCtx); ok && proc.Valid {
		if len(claims.ID) < 1 {
			return database.UserStruct{}, errors.New("id not defined")
		}

		users := database.GetUsers(db, claims.ID, sql.NullString{Valid: true, String: ""})
		if len(users) < 1 {
			return database.UserStruct{}, errors.New("user not found")
		}

		return users[0], nil
	}

	return database.UserStruct{}, err
}
