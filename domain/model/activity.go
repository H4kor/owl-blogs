package model

import "time"

type Activity struct {
	Id        string
	Type      string
	CreatedAt time.Time
	Name      string
	Content   string
	AuthorUrl string
	Raw       string
}
