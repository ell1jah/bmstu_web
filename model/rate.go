package model

const (
	Like    = true
	Dislike = false
)

type Rate bool

type RatesCnts struct {
	LikeCnt    int
	DislikeCnt int
}
