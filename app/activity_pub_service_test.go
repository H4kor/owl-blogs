package app_test

import (
	"encoding/json"
	"owl-blogs/app"
	"testing"

	vocab "github.com/go-ap/activitypub"
	"github.com/stretchr/testify/require"
)

func TestApEncoderContext(t *testing.T) {
	bytes, err := app.ApEncoder.Marshal(vocab.ObjectNew(vocab.NoteType))
	require.NoError(t, err)
	data := map[string]interface{}{}
	err = json.Unmarshal(bytes, &data)
	require.NoError(t, err)

	require.Contains(t, data, "@context")
	require.Len(t, data["@context"], 2)
	require.Contains(t, data["@context"], "https://www.w3.org/ns/activitystreams")
	require.Contains(t, data["@context"], map[string]interface{}{
		"Hashtag": "https://www.w3.org/ns/activitystreams#Hashtag",
	})
}
