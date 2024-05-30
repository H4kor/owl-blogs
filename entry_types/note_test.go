package entrytypes_test

import (
	entrytypes "owl-blogs/entry_types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoteTags(t *testing.T) {
	n := entrytypes.Note{}
	n.SetMetaData(&entrytypes.NoteMetaData{
		Content: "#tag1 hello #tagTwo #not-a-tag",
	})

	tags := n.Tags()
	require.Len(t, tags, 3)
	require.Equal(t, tags[0], "tag1")
	require.Equal(t, tags[1], "tagTwo")
	require.Equal(t, tags[2], "not")
}
