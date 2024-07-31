package models

import "time"

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"` // 这里的string表示前端可以通过string类型与后端的int64相互转换
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id,string" db:"community_id" binding:"required"`
	Status      int64     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*CommunityDetail `json:"community"`
}
