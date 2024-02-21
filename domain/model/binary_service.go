package model

type BinaryStorageInterface interface {
	Create(name string, file []byte) (*BinaryFile, error)
}
