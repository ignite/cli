package randstr

import (
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// Runes generates a random string with n length from runes.
func Runes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
