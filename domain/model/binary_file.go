package model

import (
	"mime"
	"strings"
)

type BinaryFile struct {
	Id   string
	Name string
	Data []byte
}

func (b *BinaryFile) Mime() string {
	parts := strings.Split(b.Name, ".")
	if len(parts) < 2 {
		return "application/octet-stream"
	}
	t := mime.TypeByExtension("." + parts[len(parts)-1])
	if t == "" {
		return "application/octet-stream"
	}
	return t
}
