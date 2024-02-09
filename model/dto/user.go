package dto

import (
	"github.com/ell1jah/bmstu_web/model"
)

type RespGetMe struct {
	ID    uint64 `json:"userID"`
	Login string `json:"login"`
}

func RespGetMeFromUser(user *model.User) *RespGetMe {
	return &RespGetMe{
		ID:    user.ID,
		Login: user.Login,
	}
}

type ReqСhangePass struct {
	Old string `json:"oldPassword" valid:"minstringlength(5)"`
	New string `json:"newPassword" valid:"minstringlength(5)"`
}

func (rcp *ReqСhangePass) ToChangePass() *model.UserChangePass {
	return &model.UserChangePass{
		Old: rcp.Old,
		New: rcp.New,
	}
}

type ReqSign struct {
	Login    string `json:"login" valid:"minstringlength(5)"`
	Password string `json:"password" valid:"minstringlength(5)"`
}

func (rs *ReqSign) ToUser() *model.User {
	return &model.User{
		Login:    rs.Login,
		Password: rs.Password,
	}
}

type RespToken struct {
	Token string `json:"token"`
}

func RespTokenFromString(token string) *RespToken {
	return &RespToken{
		Token: token,
	}
}
