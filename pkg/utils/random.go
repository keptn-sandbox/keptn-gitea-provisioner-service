package utils

import (
	"math/rand"
)

var /*const*/ letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"§$%&/()=?^°`,.-#+;:_'*")

// GenerateRandomString generates a random string with the length n, for generating the string a fixed alphabet is used
// that only contains printable ASCII characters
func GenerateRandomString(n int) string {
	s := make([]rune, n)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
