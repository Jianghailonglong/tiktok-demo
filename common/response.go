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

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
