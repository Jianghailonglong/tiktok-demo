package common

import (
	"mime/multipart"
)

// UserRequest 登录、注册请求内容，Param
type UserRequest struct {
	Username string `json:"username" binding:"gt=0,lte=32"`
	Password string `json:"password" binding:"gt=5,lte=32"`
}

// UserInfoRequest 用户信息请求内容，Param
type UserInfoRequest struct {
	UserId string `json:"user_id" binding:"gt=0"`
	Token  string `json:"token"`
}

// FeedRequest 视频流请求内容，Param
type FeedRequest struct {
	LatestTime string `json:"latest_time,omitempty"`
	Token      string `json:"token,omitempty"`
}

// PublishVideoRequest 发布视频请求内容，Body
type PublishVideoRequest struct {
	Data  multipart.File `json:"data"`
	Token string         `json:"token"`
	Title string         `json:"title"`
}

// VideoPublishListRequest 视频发布列表请求内容，Param
type VideoPublishListRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type ChatMessageRequest struct {
	Token      string `json:"token"`
	ToUserId   string `json:"to_user_id"`
	ActionType string `json:"action_type"`
	Content    string `json:"content"`
}

type ChatHistoryRequest struct{
	Token      string `json:"token"`
	ToUserId   string `json:"to_user_id"`
}
