package logic

import (
	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
)

type CommentRepository interface {
	GetPostComments(postId uint64) ([]*model.Comment, error)
	CreateComment(comment *model.Comment) error
}

type UserRepository interface {
	GetUserByID(id uint64) (*model.User, error)
}

type logic struct {
	commentRepository CommentRepository
	userRepository    UserRepository
}

func NewLogic(commentRepository CommentRepository, userRepository UserRepository) *logic {
	return &logic{
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (l *logic) GetPostComments(postId uint64) ([]*model.Comment, error) {
	comments, err := l.commentRepository.GetPostComments(postId)
	if err != nil {
		return nil, errors.Wrap(err, "comment repository error")
	}

	for _, comment := range comments {
		err = l.addUserInfo(comment)
		if err != nil {
			return nil, errors.Wrap(err, "addUserInfo error")
		}
	}

	return comments, nil
}

func (l *logic) CreateComment(comment *model.Comment) error {
	err := l.commentRepository.CreateComment(comment)
	if err != nil {
		return errors.Wrap(err, "comment repository error")
	}

	err = l.addUserInfo(comment)
	if err != nil {
		return errors.Wrap(err, "addUserInfo error")
	}

	return nil
}

func (l *logic) addUserInfo(comment *model.Comment) error {
	user, err := l.userRepository.GetUserByID(comment.UserID)
	if err != nil {
		return errors.Wrap(err, "user repository error")
	}

	comment.UserName = user.Login
	return nil
}
