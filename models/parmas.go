package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamCommunity struct {
	ID           int64  `json:"community_id" db:"community_id"`
	Name         string `json:"community_name" db:"community_name"`
	Introduction string `json:"introduction" db:"community_name"`
}

type ParamPostList struct {
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"`
}

type ParamVote struct {
	PostID    string `json:"post_id" binding:"required"`
	Direction int64  `json:"direction,string" binding:"oneof=1 0 -1"`
}
