package model

import "time"

type Comment struct {
	ID     uint64
	UserID uint64
	PostID uint64
	Date   time.Time
	Body   string
}
