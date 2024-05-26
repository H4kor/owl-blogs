package app

import (
	"owl-blogs/domain/model"

	vocab "github.com/go-ap/activitypub"
)

// this file contains interfaces which CAN be implemented by entry types
// to support additional distribution channels.

// ToActivityPub can be implemented to send entries to followers via the ActivityPub protocol.
// Following fields will be filled in by the ActivityPubService:
//   - ID
//   - AttributedTo
//   - To
//
// The site configuration and binary service are provided
type ToActivityPub interface {
	ActivityObject(siteCfg model.SiteConfig, binSvc BinaryService) vocab.Object
}
