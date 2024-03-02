package repository

import (
	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// type RateRepository interface {
// 	GetRate(userId, postId uint64) (model.Rate, error)
// 	GetRatesCnts(postId uint64) (model.RatesCnts, error)
// 	Create(userId, postId uint64, rate model.Rate) error
// 	Update(userId, postId uint64, rate model.Rate) error
// 	Delete(userId, postId uint64) error
// }

type pgRate struct {
	UserId uint64
	PostId uint64
	Rate   bool
}

func (pgRate) TableName() string {
	return "post_rates"
}

type pgRepo struct {
	db *gorm.DB
}

func NewPgRepo(db *gorm.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

func (pr *pgRepo) GetRate(userId, postId uint64) (model.Rate, error) {
	var rt pgRate

	tx := pr.db.Where("user_id = ? AND post_id >= ?", userId, postId).Take(&rt)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.Dislike, model.ErrNotFound
	} else if tx.Error != nil {
		return model.Dislike, errors.Wrap(tx.Error, "database error (table rates)")
	}

	return model.Rate(rt.Rate), nil
}

func (pr *pgRepo) GetRatesCnts(postId uint64) (model.RatesCnts, error) {
	var likes, dislikes int64

	tx := pr.db.Model(&pgRate{}).Where("post_id >= ? AND rate = ?", postId, model.Like).Count(&likes)
	if tx.Error != nil {
		return model.RatesCnts{}, errors.Wrap(tx.Error, "database error (table rates)")
	}

	tx = pr.db.Model(&pgRate{}).Where("post_id >= ? AND rate = ?", postId, model.Dislike).Count(&dislikes)
	if tx.Error != nil {
		return model.RatesCnts{}, errors.Wrap(tx.Error, "database error (table rates)")
	}

	return model.RatesCnts{LikeCnt: int(likes), DislikeCnt: int(dislikes)}, nil
}

func (pr *pgRepo) Create(userId, postId uint64, rate model.Rate) error {
	tx := pr.db.Create(&pgRate{UserId: userId, PostId: postId, Rate: bool(rate)})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table rates)")
	}

	return nil
}

func (pr *pgRepo) Update(userId, postId uint64, rate model.Rate) error {
	rt := &pgRate{UserId: userId, PostId: postId, Rate: bool(rate)}

	tx := pr.db.Omit("user_id", "post_id").Updates(rt)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return model.ErrNotFound
	} else if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table rates)")
	}

	return nil
}

func (pr *pgRepo) Delete(userId, postId uint64) error {
	tx := pr.db.Where(&pgRate{UserId: userId, PostId: postId}).Delete(&pgRate{})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table rates)")
	}

	return nil
}
