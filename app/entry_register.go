package app

import (
	"errors"
	"owl-blogs/domain/model"
	"reflect"
)

type EntryTypeRegistry struct {
	types map[string]model.Entry
}

func NewEntryTypeRegistry() *EntryTypeRegistry {
	return &EntryTypeRegistry{types: map[string]model.Entry{}}
}

func (r *EntryTypeRegistry) entryType(entry model.Entry) string {
	return reflect.TypeOf(entry).Elem().Name()
}

func (r *EntryTypeRegistry) Register(entry model.Entry) error {
	t := r.entryType(entry)
	if _, ok := r.types[t]; ok {
		return errors.New("entry type already registered")
	}
	r.types[t] = entry
	return nil
}

func (r *EntryTypeRegistry) Types() []model.Entry {
	types := []model.Entry{}
	for _, t := range r.types {
		types = append(types, t)
	}
	return types
}

func (r *EntryTypeRegistry) TypeName(entry model.Entry) (string, error) {
	t := r.entryType(entry)
	if _, ok := r.types[t]; !ok {
		return "", errors.New("entry type not registered")
	}
	return t, nil
}

func (r *EntryTypeRegistry) Type(name string) (model.Entry, error) {
	if _, ok := r.types[name]; !ok {
		return nil, errors.New("entry type not registered")
	}
	return r.types[name], nil
}
