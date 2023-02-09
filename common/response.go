package common

// UserResponse 登录、注册响应内容
type UserResponse struct {
	Response
	UserId int32  `json:"user_id"`
	Token  string `json:"token"`
}

// UserInfoResponse 显示用户信息响应内容
type UserInfoResponse struct {
	Response
	UserInfo User `json:"user"`
}

// FeedResponse 视频流响应内容
type FeedResponse struct {
	Response
	NextTime  int64   `json:"next_time"` // 作为下次请求的latest_time
	VideoList []Video `json:"video_list"`
}

// VideoPublishListResponse 视频发布列表响应内容
type VideoPublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// CommentListResponse 评论表的响应结构体
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentActionResponse 评论操作的相应结构体
type CommentActionResponse struct {
	Response
	Comment *Comment `json:"comment,omitempty"`
}

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// FollowerListResponse 获取粉丝列表响应内容
type FollowerListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// FollowListResponse 获取关注列表响应内容
type FollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
