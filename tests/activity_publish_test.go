package tests

import (
	entrytypes "owl-blogs/entry_types"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestEntryIsSent(t *testing.T) {
	//setup
	app := DefaultTestApp()
	srv := adaptor.FiberApp(app.FiberApp)
	mock := NewMockAPServer()
	defer mock.Server.Close()

	EnsureFollowed(t, srv, mock, mock.MockActorUrl("1"))
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 1, len(mock.Retrieved["1"]))

	note := entrytypes.Note{}
	note.SetMetaData(&entrytypes.NoteMetaData{
		Content: "test note",
	})
	now := time.Now()
	note.SetPublishedAt(&now)

	app.EntryService.Create(&note)
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 2, len(mock.Retrieved["1"]))

	msg := mock.Retrieved["1"][1]

	require.Equal(t, "Create", msg["type"])
	require.Contains(t, msg, "id")
	require.Contains(t, msg, "published")
	require.Contains(t, msg["object"].(map[string]interface{})["content"], "test note")

}

func TestEntrySentToLinks(t *testing.T) {
	//setup
	app := DefaultTestApp()
	adaptor.FiberApp(app.FiberApp)
	mock := NewMockAPServer()
	defer mock.Server.Close()

	actor := mock.MockActorUrl("1")
	activity := mock.MockActivityUrl(actor, "1")

	note := entrytypes.Note{}
	note.SetMetaData(&entrytypes.NoteMetaData{
		Content: "test note <" + activity + ">",
	})
	now := time.Now()
	note.SetPublishedAt(&now)

	app.EntryService.Create(&note)
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 1, len(mock.Retrieved["1"]))

	msg := mock.Retrieved["1"][0]

	require.Equal(t, "Create", msg["type"])
	require.Contains(t, msg, "id")
	require.Contains(t, msg, "published")
	require.Contains(t, msg["object"].(map[string]interface{})["content"], "test note")

}

func TestEntrySentToLinkedActor(t *testing.T) {
	//setup
	app := DefaultTestApp()
	adaptor.FiberApp(app.FiberApp)
	mock := NewMockAPServer()
	defer mock.Server.Close()

	actor := mock.MockActorUrl("1")

	note := entrytypes.Note{}
	note.SetMetaData(&entrytypes.NoteMetaData{
		Content: "test note <" + actor + ">",
	})
	now := time.Now()
	note.SetPublishedAt(&now)

	app.EntryService.Create(&note)
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 1, len(mock.Retrieved["1"]))

	msg := mock.Retrieved["1"][0]

	require.Equal(t, "Create", msg["type"])
	require.Contains(t, msg, "id")
	require.Contains(t, msg, "published")
	require.Contains(t, msg["object"].(map[string]interface{})["content"], "test note")

}

func TestEntryAutoAttachments(t *testing.T) {
	//setup
	app := DefaultTestApp()
	srv := adaptor.FiberApp(app.FiberApp)
	mock := NewMockAPServer()
	defer mock.Server.Close()

	EnsureFollowed(t, srv, mock, mock.MockActorUrl("1"))
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 1, len(mock.Retrieved["1"]))

	bin, err := app.BinaryService.Create("image.png", []byte("fooo"))
	require.NoError(t, err)

	note := entrytypes.Note{}
	note.SetMetaData(&entrytypes.NoteMetaData{
		Content: "![This is the alt text](/media/" + bin.Id + ")",
	})
	now := time.Now()
	note.SetPublishedAt(&now)

	app.EntryService.Create(&note)
	time.Sleep(50 * time.Millisecond)

	require.Equal(t, 1, len(mock.Retrieved))
	require.Equal(t, 2, len(mock.Retrieved["1"]))

	msg := mock.Retrieved["1"][1]

	require.Equal(t, "Create", msg["type"])
	require.Contains(t, msg, "id")
	require.Contains(t, msg, "published")
	require.Equal(t, msg["object"].(map[string]interface{})["attachment"].(map[string]interface{})["name"], "This is the alt text")
	require.Contains(t, msg["object"].(map[string]interface{})["attachment"].(map[string]interface{})["url"], "image.png")

}
