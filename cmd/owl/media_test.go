package main

import (
	"net/http/httptest"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMediaWithSpace(t *testing.T) {
	db := test.NewMockDb()
	owlApp := App(db)
	app := owlApp.FiberApp

	_, err := owlApp.BinaryService.Create("name with space.jpg", []byte("111"))
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/media/name%20with%20space.jpg", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

}
