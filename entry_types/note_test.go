package entrytypes_test

import (
	entrytypes "owl-blogs/entry_types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoteTags(t *testing.T) {
	n := entrytypes.Note{}
	n.SetMetaData(&entrytypes.NoteMetaData{
		Content: "#tag1 hello #tagTwo #a-tag",
	})

	tags := n.Tags()
	require.Len(t, tags, 3)
	require.Contains(t, tags, "tag1")
	require.Contains(t, tags, "tagTwo")
	require.Contains(t, tags, "a-tag")
}
