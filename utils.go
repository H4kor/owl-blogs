package owl

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[k.Int64()]
	}
	return string(b)
}
