package logic

import (
	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserByID(id uint64) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
}

type logic struct {
	userRepository UserRepository
}

func NewLogic(userRepository UserRepository) *logic {
	return &logic{
		userRepository: userRepository,
	}
}

func (l *logic) GetUserByID(id uint64) (*model.User, error) {
	user, err := l.userRepository.GetUserByID(id)

	if err != nil {
		return nil, errors.Wrap(err, "user repository error")
	}

	return user, nil
}

func (l *logic) ChangePass(chpass *model.UserChangePass) error {
	if chpass.Old == chpass.New {
		return model.ErrConflictPassword
	}

	user, err := l.userRepository.GetUserByID(chpass.ID)
	if err != nil {
		return errors.Wrap(err, "user repository error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(chpass.Old))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return model.ErrInvalidPassword
	} else if err != nil {
		return errors.Wrap(err, "bcrypt error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(chpass.New), 8)
	if err != nil {
		return errors.Wrap(err, "bcrypt error")
	}

	user.Password = string(hashedPassword)
	_, err = l.userRepository.UpdateUser(user)
	if err != nil {
		return errors.Wrap(err, "user repository error")
	}

	return nil
}

func (l *logic) SignIn(user *model.User) (*model.User, error) {
	gotUser, err := l.userRepository.GetUserByLogin(user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "user repository error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(gotUser.Password), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, model.ErrInvalidPassword
	} else if err != nil {
		return nil, errors.Wrap(err, "bcrypt error")
	}

	gotUser.Password = ""
	return gotUser, nil
}

func (l *logic) SignUp(user *model.User) (*model.User, error) {
	_, err := l.userRepository.GetUserByLogin(user.Login)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrap(err, "user repository error")
	} else if err == nil {
		return nil, model.ErrConflictNickname
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, errors.Wrap(err, "bcrypt error")
	}

	user.Password = string(hashedPassword)
	user, err = l.userRepository.CreateUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "user repository error")
	}

	return user, nil
}
