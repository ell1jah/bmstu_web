package model

import "github.com/pkg/errors"

var (
	ErrNotFound            = errors.New("item is not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrConflictPassword    = errors.New("new password is the same as old")
	ErrConflictNickname    = errors.New("nickname already exists")
	ErrConflictEmail       = errors.New("email already exists")
	ErrBadRequest          = errors.New("bad request")
	ErrConflictFriend      = errors.New("friend already exists")
	ErrUnauthorized        = errors.New("no cookie")
	ErrInternalServerError = errors.New("internal server error")
	ErrEmptyCsrf           = errors.New("empty csrf token")
	ErrInvalidCsrf         = errors.New("invalid csrf")
	ErrPermissionDenied    = errors.New("permission denied")
)
