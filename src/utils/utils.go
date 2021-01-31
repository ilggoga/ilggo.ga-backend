package utils

import (
	"crypto/rand"
	"math/big"
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
