package app

import (
	"errors"
	"reflect"
)

type TypeRegistry[T any] struct {
	types map[string]T
}

func NewTypeRegistry[T any]() *TypeRegistry[T] {
	return &TypeRegistry[T]{types: map[string]T{}}
}

func (r *TypeRegistry[T]) entryType(entry T) string {
	return reflect.TypeOf(entry).Elem().Name()
}

func (r *TypeRegistry[T]) Register(entry T) error {
	t := r.entryType(entry)
	if _, ok := r.types[t]; ok {
		return errors.New("entry type already registered")
	}
	r.types[t] = entry
	return nil
}

func (r *TypeRegistry[T]) Types() []T {
	types := []T{}
	for _, t := range r.types {
		types = append(types, t)
	}
	return types
}

func (r *TypeRegistry[T]) TypeName(entry T) (string, error) {
	t := r.entryType(entry)
	if _, ok := r.types[t]; !ok {
		return "", errors.New("entry type not registered")
	}
	return t, nil
}

func (r *TypeRegistry[T]) Type(name string) (T, error) {
	if _, ok := r.types[name]; !ok {
		return *new(T), errors.New("entry type not registered")
	}

	val := reflect.ValueOf(r.types[name])
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	newEntry := reflect.New(val.Type()).Interface().(T)

	return newEntry, nil
}
