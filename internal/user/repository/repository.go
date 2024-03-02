package repository

import (
	"github.com/ell1jah/bmstu_web/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type pgUser struct {
	ID       uint64
	Login    string
	Password string
}

func (u pgUser) toModelUser() *model.User {
	return &model.User{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func fromModelUser(u *model.User) *pgUser {
	return &pgUser{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func (pgUser) TableName() string {
	return "users"
}

type pgRepo struct {
	db *gorm.DB
}

func NewPgRepo(db *gorm.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

func (pr *pgRepo) GetUserByID(id uint64) (*model.User, error) {
	var usr pgUser

	tx := pr.db.Where("id = ?", id).Take(&usr)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return usr.toModelUser(), nil
}

func (pr *pgRepo) GetUserByLogin(login string) (*model.User, error) {
	var usr pgUser

	tx := pr.db.Where("login = ?", login).Take(&usr)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return usr.toModelUser(), nil
}

func (pr *pgRepo) UpdateUser(user *model.User) (*model.User, error) {
	oldUser := fromModelUser(user)

	tx := pr.db.Omit("id").Updates(oldUser)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return user, nil
}

func (pr *pgRepo) CreateUser(user *model.User) (*model.User, error) {
	pgUsr := fromModelUser(user)

	tx := pr.db.Create(pgUsr)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, model.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	user.ID = pgUsr.ID

	return user, nil
}
