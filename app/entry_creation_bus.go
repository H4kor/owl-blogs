package app

import "owl-blogs/domain/model"

type Subscriber interface{}
type EntryCreatedSubscriber interface {
	NotifyEntryCreated(entry model.Entry)
}

type EntryUpdatedSubscriber interface {
	NotifyEntryUpdated(entry model.Entry)
}
type EntryDeletedSubscriber interface {
	NotifyEntryDeleted(entry model.Entry)
}

type EventBus struct {
	subscribers []Subscriber
}

func NewEntryCreationBus() *EventBus {
	return &EventBus{subscribers: make([]Subscriber, 0)}
}

func (b *EventBus) Subscribe(subscriber Subscriber) {
	b.subscribers = append(b.subscribers, subscriber)
}

func (b *EventBus) NotifyCreated(entry model.Entry) {
	for _, subscriber := range b.subscribers {
		if sub, ok := subscriber.(EntryCreatedSubscriber); ok {
			go sub.NotifyEntryCreated(entry)
		}
	}
}

func (b *EventBus) NotifyUpdated(entry model.Entry) {
	for _, subscriber := range b.subscribers {
		if sub, ok := subscriber.(EntryUpdatedSubscriber); ok {
			go sub.NotifyEntryUpdated(entry)
		}
	}
}

func (b *EventBus) NotifyDeleted(entry model.Entry) {
	for _, subscriber := range b.subscribers {
		if sub, ok := subscriber.(EntryDeletedSubscriber); ok {
			go sub.NotifyEntryDeleted(entry)
		}
	}
}
