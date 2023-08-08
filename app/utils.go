package app

import (
	"math/rand"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func UrlToEntryId(url string) string {
	parts := strings.Split(url, "/")
	if parts[len(parts)-1] == "" {
		return parts[len(parts)-2]
	} else {
		return parts[len(parts)-1]
	}
}
