package model

import "time"

type Post struct {
	ID          uint64
	UserID      uint64
	Date        time.Time
	ImageID     uint64
	Category    string
	Brand       string
	Description string
	Link        string
}
