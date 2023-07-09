package model_test

import (
	"owl-blogs/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMimeType(t *testing.T) {
	bin := model.BinaryFile{Name: "test.jpg"}
	require.Equal(t, "image/jpeg", bin.Mime())
}
