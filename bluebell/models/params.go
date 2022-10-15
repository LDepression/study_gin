package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//ParamsSignUp 用于来存放登录的相关的参数
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
}

type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ParamsVoteData struct {
	//UserID 从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`                       //帖子id
	Direction int    `json:"direction,string" binding:"required,oneof=1 0 -1"` //赞成票(1)还是反对票(-1)取消投票(0)
}

//ParamPostList 获取帖子列表的querystring参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

type ParamsCommunityPostList struct {
	ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
