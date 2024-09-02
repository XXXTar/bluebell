package models

//定义请求的参数

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//注册请求参数
type ParamSignup struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

//登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//投票数据
type ParamVoteData struct {
	//UserID 从当前请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`               //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1" ` //赞成(1)反对(-1)票
}

//获取帖子列表query string参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   //可以为空
	Page        int64  `json:"page" form:"page" example:"1"`       //页码
	Size        int64  `json:"size" form:"size" example:"10"`      //每页数据量
	Order       string `json:"order" form:"order" example:"score"` //排序依据
}
