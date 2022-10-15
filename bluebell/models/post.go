package models

import "time"

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	State       int32     `json:"state" db:"state"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}
type ApiPostDetails struct {
	AuthorName       string `json:"author_name"`
	VoteNum          int64  `json:"vote_num"`
	*Post            `json:"post"`
	*CommunityDetail `json:"community"`
}
