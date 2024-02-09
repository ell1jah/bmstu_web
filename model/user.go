package model

type User struct {
	ID       uint64
	Login    string
	Password string
}

type UserChangePass struct {
	ID  uint64
	Old string
	New string
}
