package repository

import (
	"time"

	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pgPost struct {
	ID          uint64
	UserID      uint64
	CreatedAt   time.Time
	ImageID     string
	Category    string
	Sex         string
	Brand       string
	Description string
	Link        string
}

func (p pgPost) toModelPost() *model.Post {
	return &model.Post{
		ID:          p.ID,
		UserID:      p.UserID,
		Date:        p.CreatedAt,
		ImageID:     p.ImageID,
		Category:    p.Category,
		Sex:         p.Sex,
		Brand:       p.Brand,
		Description: p.Description,
		Link:        p.Link,
	}
}

func toModelPosts(pg []*pgPost) []*model.Post {
	posts := make([]*model.Post, len(pg))

	for i := range posts {
		posts[i] = pg[i].toModelPost()
	}

	return posts
}

func fromModelPost(p *model.Post) *pgPost {
	return &pgPost{
		ID:          p.ID,
		UserID:      p.UserID,
		CreatedAt:   p.Date,
		ImageID:     p.ImageID,
		Category:    p.Category,
		Sex:         p.Sex,
		Brand:       p.Brand,
		Description: p.Description,
		Link:        p.Link,
	}
}

func (pgPost) TableName() string {
	return "posts"
}

type pgRepo struct {
	db *gorm.DB
}

func NewPgRepo(db *gorm.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

func (pr *pgRepo) GetPost(postId uint64) (*model.Post, error) {
	var pst pgPost

	tx := pr.db.Where("id = ?", postId).Take(&pst)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return pst.toModelPost(), nil
}

func (pr *pgRepo) GetUsersPosts(ownerId uint64) ([]*model.Post, error) {
	posts := make([]*pgPost, 0, 10)

	tx := pr.db.Where(&pgPost{UserID: ownerId}).Order("id desc").Find(&posts)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return toModelPosts(posts), nil
}

func (pr *pgRepo) GetPostsWithParams(params model.PostParams) ([]*model.Post, error) {
	posts := make([]*pgPost, 0, 10)

	tx := pr.db.Where(fromModelPost(params.ToPost())).Order("id desc").Find(&posts)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return toModelPosts(posts), nil
}

func (pr *pgRepo) CreatePost(post *model.Post) error {
	post.Date = time.Now()
	pgPost := fromModelPost(post)

	tx := pr.db.Create(pgPost)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	post.ID = pgPost.ID
	return nil
}

func (pr *pgRepo) DeletePost(postId uint64) error {
	tx := pr.db.Delete(&pgPost{}, postId)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}
