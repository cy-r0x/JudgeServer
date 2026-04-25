package usercsv

import (
	"crypto/rand"
	"math/big"
)

func generatePassword() string {
	const length = 8
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@$()!"

	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range password {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "ChangeMe123!"
		}
		password[i] = charset[num.Int64()]
	}

	return string(password)
}
