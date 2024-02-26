package model

import "time"

type Post struct {
	ID          uint64
	UserID      uint64
	UserName    string
	Date        time.Time
	ImageID     string
	Category    string
	Sex         string
	Brand       string
	Description string
	Link        string
	LikeCnt     int
	DislikeCnt  int
	IsLiked     bool
	IsDisliked  bool
}

type PostParams struct {
	Sex      string
	Category string
}

func (pp PostParams) ToPost() *Post {
	return &Post{
		Sex:      pp.Sex,
		Category: pp.Category,
	}
}
