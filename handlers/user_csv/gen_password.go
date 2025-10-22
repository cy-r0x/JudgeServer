package usercsv

import (
	"math/rand"
	"strings"
)

func generatePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#&$"
	const length = 8
	var password strings.Builder
	for range length {
		password.WriteByte(charset[rand.Intn(len(charset))])
	}
	return password.String()
}
