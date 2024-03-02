package dto

import (
	"time"

	"github.com/ell1jah/bmstu_web/model"
)

type RespComment struct {
	ID       uint64    `json:"commentID"`
	UserID   uint64    `json:"creatorID"`
	UserName string    `json:"creatorName"`
	PostID   uint64    `json:"postID"`
	Date     time.Time `json:"createDate"`
	Body     string    `json:"message"`
}

func RespCommentFromComment(c *model.Comment) *RespComment {
	return &RespComment{
		ID:       c.ID,
		UserID:   c.UserID,
		UserName: c.UserName,
		PostID:   c.PostID,
		Date:     c.Date,
		Body:     c.Body,
	}
}

func RespCommentsFromComments(comments []*model.Comment) []*RespComment {
	resp := make([]*RespComment, len(comments))
	for i := range resp {
		resp[i] = RespCommentFromComment(comments[i])
	}

	return resp
}

type ReqComment struct {
	Body string `json:"message" valid:"minstringlength(1)"`
}

func (rc *ReqComment) ToComment() *model.Comment {
	return &model.Comment{
		Body: rc.Body,
	}
}
