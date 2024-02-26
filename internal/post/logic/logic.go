package logic

import (
	"os"

	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
)

const (
	imageDir = "images/"
	pngExt   = ".png"
)

type PostRepository interface {
	GetPost(postId uint64) (*model.Post, error)
	GetUsersPosts(ownerId uint64) ([]*model.Post, error)
	GetPostsWithParams(params model.PostParams) ([]*model.Post, error)
	CreatePost(post *model.Post) error
	DeletePost(postId uint64) error
}

type UserRepository interface {
	GetUserByID(id uint64) (*model.User, error)
}

type RateRepository interface {
	GetRate(userId, postId uint64) (model.Rate, error)
	GetRatesCnts(postId uint64) (model.RatesCnts, error)
	Create(userId, postId uint64, rate model.Rate) error
	Update(userId, postId uint64, rate model.Rate) error
	Delete(userId, postId uint64) error
}

type logic struct {
	postRepository PostRepository
	userRepository UserRepository
	rateRepository RateRepository
}

func NewLogic(postRepository PostRepository, userRepository UserRepository, rateRepository RateRepository) *logic {
	return &logic{
		postRepository: postRepository,
		userRepository: userRepository,
		rateRepository: rateRepository,
	}
}

func (l *logic) GetPost(userId, postId uint64) (*model.Post, error) {
	post, err := l.postRepository.GetPost(postId)
	if err != nil {
		return nil, errors.Wrap(err, "post repository error")
	}

	err = l.addUserInfo(post)
	if err != nil {
		return nil, errors.Wrap(err, "addUserInfo error")
	}

	err = l.addRateInfo(userId, post)
	if err != nil {
		return nil, errors.Wrap(err, "addRateInfo error")
	}

	return post, nil
}

func (l *logic) GetUsersPosts(askerId, ownerId uint64) ([]*model.Post, error) {
	posts, err := l.postRepository.GetUsersPosts(ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "post repository error")
	}

	for _, post := range posts {
		err = l.addUserInfo(post)
		if err != nil {
			return nil, errors.Wrap(err, "addUserInfo error")
		}

		err = l.addRateInfo(askerId, post)
		if err != nil {
			return nil, errors.Wrap(err, "addRateInfo error")
		}
	}

	return posts, nil
}

func (l *logic) GetPostsWithParams(userId uint64, params model.PostParams) ([]*model.Post, error) {
	posts, err := l.postRepository.GetPostsWithParams(params)
	if err != nil {
		return nil, errors.Wrap(err, "post repository error")
	}

	for _, post := range posts {
		err = l.addUserInfo(post)
		if err != nil {
			return nil, errors.Wrap(err, "addUserInfo error")
		}

		err = l.addRateInfo(userId, post)
		if err != nil {
			return nil, errors.Wrap(err, "addRateInfo error")
		}
	}

	return posts, nil
}

func (l *logic) CreatePost(post *model.Post) error {
	if _, err := os.Stat(imageDir + post.ImageID + pngExt); errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(model.ErrBadRequest, "no image")
	} else if err != nil {
		return errors.Wrap(err, "can't find image")
	}

	err := l.postRepository.CreatePost(post)
	if err != nil {
		return errors.Wrap(err, "post repository error")
	}

	err = l.addUserInfo(post)
	if err != nil {
		return errors.Wrap(err, "addUserInfo error")
	}

	return nil
}

func (l *logic) DeletePost(userId, postId uint64) error {
	post, err := l.postRepository.GetPost(postId)
	if err != nil {
		return errors.Wrap(err, "post repository error")
	}

	if post.UserID != userId {
		return model.ErrPermissionDenied
	}

	err = l.postRepository.DeletePost(postId)
	if err != nil {
		return errors.Wrap(err, "post repository error")
	}

	return nil
}

func (l *logic) LikePost(userId, postId uint64) error {
	rate, err := l.rateRepository.GetRate(userId, postId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return errors.Wrap(err, "rate repository error")
	} else if err != nil {
		err = l.rateRepository.Create(userId, postId, model.Like)
		if err != nil {
			return errors.Wrap(err, "rate repository error")
		}
	} else {
		if rate != model.Like {
			err = l.rateRepository.Update(userId, postId, model.Like)
			if err != nil {
				return errors.Wrap(err, "rate repository error")
			}
		}
	}

	return nil
}

func (l *logic) DislikePost(userId, postId uint64) error {
	rate, err := l.rateRepository.GetRate(userId, postId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return errors.Wrap(err, "rate repository error")
	} else if err != nil {
		err = l.rateRepository.Create(userId, postId, model.Dislike)
		if err != nil {
			return errors.Wrap(err, "rate repository error")
		}
	} else {
		if rate != model.Dislike {
			err = l.rateRepository.Update(userId, postId, model.Dislike)
			if err != nil {
				return errors.Wrap(err, "rate repository error")
			}
		}
	}

	return nil
}

func (l *logic) UnratePost(userId, postId uint64) error {
	err := l.rateRepository.Delete(userId, postId)
	if err != nil {
		return errors.Wrap(err, "rate repository error")
	}

	return nil
}

func (l *logic) addUserInfo(post *model.Post) error {
	user, err := l.userRepository.GetUserByID(post.UserID)
	if err != nil {
		return errors.Wrap(err, "user repository error")
	}

	post.UserName = user.Login
	return nil
}

func (l *logic) addRateInfo(userId uint64, post *model.Post) error {
	rate, err := l.rateRepository.GetRate(userId, post.ID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return errors.Wrap(err, "rate repository error")
	} else if err != nil {
		post.IsLiked = false
		post.IsDisliked = false
	} else {
		if rate == model.Like {
			post.IsLiked = true
			post.IsDisliked = false
		} else {
			post.IsLiked = false
			post.IsDisliked = true
		}
	}

	rateCnt, err := l.rateRepository.GetRatesCnts(post.ID)
	if err != nil {
		return errors.Wrap(err, "rate repository error")
	}

	post.LikeCnt = rateCnt.LikeCnt
	post.DislikeCnt = rateCnt.DislikeCnt
	return nil
}
