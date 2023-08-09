package app

import (
	"owl-blogs/domain/model"
)

type EntryTypeRegistry = TypeRegistry[model.Entry]

func NewEntryTypeRegistry() *EntryTypeRegistry {
	return NewTypeRegistry[model.Entry]()
}
