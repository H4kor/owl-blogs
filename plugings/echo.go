package plugings

import (
	"fmt"
	"owl-blogs/app"
	"owl-blogs/domain/model"
)

type Echo struct {
}

func NewEcho(bus *app.EventBus) *Echo {
	echo := &Echo{}
	bus.Subscribe(echo)
	return echo
}

func (e *Echo) NotifyEntryCreated(entry model.Entry) {
	fmt.Println("Entry Created:")
	fmt.Println("\tID: ", entry.ID())
	fmt.Println("\tTitle: ", entry.Title())
	fmt.Println("\tPublishedAt: ", entry.PublishedAt())
}
