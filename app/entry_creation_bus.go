package app

import "owl-blogs/domain/model"

type EntryCreationSubscriber interface {
	NotifyEntryCreation(entry model.Entry)
}

type EntryCreationBus struct {
	subscribers []EntryCreationSubscriber
}

func NewEntryCreationBus() *EntryCreationBus {
	return &EntryCreationBus{subscribers: make([]EntryCreationSubscriber, 0)}
}

func (b *EntryCreationBus) Subscribe(subscriber EntryCreationSubscriber) {
	b.subscribers = append(b.subscribers, subscriber)
}

func (b *EntryCreationBus) Notify(entry model.Entry) {
	for _, subscriber := range b.subscribers {
		subscriber.NotifyEntryCreation(entry)
	}
}
