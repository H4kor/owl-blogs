package model

// Thumbnail is a smaller version of an image BinaryFile used to save bandwidth
type Thumbnail struct {
	Id string
	// BinaryFileId is the Id of the BinaryFile from which the Thumbnail was created
	BinaryFileId string
	// Data is the raw data of the image
	Data []byte
	// MimeType of the Thumbnail, e.g. "image/jpeg"
	MimeType string
}
