package app_test

import (
	"owl-blogs/app"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryTypeNameNotExisting(t *testing.T) {
	register := app.NewEntryTypeRegistry()
	_, err := register.TypeName(&test.MockEntry{})
	require.Error(t, err)
}

func TestRegistryTypeName(t *testing.T) {
	register := app.NewEntryTypeRegistry()
	register.Register(&test.MockEntry{})
	name, err := register.TypeName(&test.MockEntry{})
	require.NoError(t, err)
	require.Equal(t, "MockEntry", name)
}
