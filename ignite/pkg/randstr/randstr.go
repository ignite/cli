package randstr

import (
	"crypto/rand"
	"math/big"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// Runes generates a random string with n length from runes.
func Runes(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		b[i] = letterRunes[num.Int64()]
	}
	return string(b)
}
