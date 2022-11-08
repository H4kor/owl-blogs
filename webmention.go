package owl

import (
	"time"
)

type WebmentionIn struct {
	Source         string    `yaml:"source"`
	Title          string    `yaml:"title"`
	ApprovalStatus string    `yaml:"approval_status"`
	RetrievedAt    time.Time `yaml:"retrieved_at"`
}

func (webmention *WebmentionIn) UpdateWith(update WebmentionIn) {
	if update.Title != "" {
		webmention.Title = update.Title
	}
	if update.ApprovalStatus != "" {
		webmention.ApprovalStatus = update.ApprovalStatus
	}
	if !update.RetrievedAt.IsZero() {
		webmention.RetrievedAt = update.RetrievedAt
	}
}

type WebmentionOut struct {
	Target     string    `yaml:"target"`
	Supported  bool      `yaml:"supported"`
	ScannedAt  time.Time `yaml:"scanned_at"`
	LastSentAt time.Time `yaml:"last_sent_at"`
}

func (webmention *WebmentionOut) UpdateWith(update WebmentionOut) {
	if update.Supported {
		webmention.Supported = update.Supported
	}
	if !update.ScannedAt.IsZero() {
		webmention.ScannedAt = update.ScannedAt
	}
	if !update.LastSentAt.IsZero() {
		webmention.LastSentAt = update.LastSentAt
	}
}
