package app

import (
	"owl-blogs/domain/model"
)

type InteractionTypeRegistry = TypeRegistry[model.Interaction]

func NewInteractionTypeRegistry() *InteractionTypeRegistry {
	return NewTypeRegistry[model.Interaction]()
}
