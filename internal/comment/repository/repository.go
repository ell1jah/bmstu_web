package repository

import (
	"time"

	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pgComment struct {
	ID        uint64
	UserID    uint64
	PostID    uint64
	CreatedAt time.Time
	Text      string
}

func (c pgComment) toModelComment() *model.Comment {
	return &model.Comment{
		ID:     c.ID,
		UserID: c.UserID,
		PostID: c.PostID,
		Date:   c.CreatedAt,
		Body:   c.Text,
	}
}

func toModelComments(pg []*pgComment) []*model.Comment {
	comments := make([]*model.Comment, len(pg))

	for i := range comments {
		comments[i] = pg[i].toModelComment()
	}

	return comments
}

func fromModelComment(c *model.Comment) *pgComment {
	return &pgComment{
		ID:        c.ID,
		UserID:    c.UserID,
		PostID:    c.PostID,
		CreatedAt: c.Date,
		Text:      c.Body,
	}
}

func (pgComment) TableName() string {
	return "comments"
}

type pgRepo struct {
	db *gorm.DB
}

func NewPgRepo(db *gorm.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

func (pr *pgRepo) GetPostComments(postId uint64) ([]*model.Comment, error) {
	comments := make([]*pgComment, 0, 10)

	tx := pr.db.Where(&pgComment{PostID: postId}).Order("id desc").Find(&comments)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table comments)")
	}

	return toModelComments(comments), nil
}

func (pr *pgRepo) CreateComment(comment *model.Comment) error {
	comment.Date = time.Now()
	pgComment := fromModelComment(comment)

	tx := pr.db.Create(pgComment)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table comments)")
	}

	comment.ID = pgComment.ID
	return nil
}
