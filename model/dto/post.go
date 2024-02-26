package dto

import (
	"time"

	"github.com/ell1jah/bmstu_web/model"
)

type RespPost struct {
	ID          uint64    `json:"postID"`
	UserID      uint64    `json:"creatorID"`
	UserName    string    `json:"creatorName"`
	Date        time.Time `json:"createDate"`
	ImageID     string    `json:"photoID"`
	Category    string    `json:"category"`
	Sex         string    `json:"sex"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	LikeCnt     int       `json:"likeCnt"`
	DislikeCnt  int       `json:"dislikeCnt"`
	IsLiked     bool      `json:"isLiked"`
	IsDisliked  bool      `json:"isDisliked"`
}

func RespPostFromPost(post *model.Post) *RespPost {
	return &RespPost{
		ID:          post.ID,
		UserID:      post.UserID,
		UserName:    post.UserName,
		Date:        post.Date,
		ImageID:     post.ImageID,
		Category:    post.Category,
		Sex:         post.Sex,
		Brand:       post.Brand,
		Description: post.Description,
		Link:        post.Link,
		LikeCnt:     post.LikeCnt,
		DislikeCnt:  post.DislikeCnt,
		IsLiked:     post.IsLiked,
		IsDisliked:  post.IsDisliked,
	}
}

func RespPostsFromPosts(posts []*model.Post) []*RespPost {
	resp := make([]*RespPost, len(posts))
	for i := range resp {
		resp[i] = RespPostFromPost(posts[i])
	}

	return resp
}

type ReqPost struct {
	ImageID     string `json:"photoID" valid:"-"`
	Category    string `json:"category" valid:"in(shoes|outerwear|underwear|accessories)"`
	Sex         string `json:"sex" valid:"in(male|female)"`
	Brand       string `json:"brand" valid:"-"`
	Description string `json:"description" valid:"-"`
	Link        string `json:"link" valid:"-"`
}

func (rp *ReqPost) ToPost() *model.Post {
	return &model.Post{
		ImageID:     rp.ImageID,
		Category:    rp.Category,
		Sex:         rp.Sex,
		Brand:       rp.Brand,
		Description: rp.Description,
		Link:        rp.Link,
	}
}

type ReqPostParams struct {
	Category string `query:"category" valid:"in(shoes|outerwear|underwear|accessories),optional"`
	Sex      string `query:"sex" valid:"in(male|female),optional"`
}

func (rpp *ReqPostParams) ToPostParams() *model.PostParams {
	return &model.PostParams{
		Category: rpp.Category,
		Sex:      rpp.Sex,
	}
}
