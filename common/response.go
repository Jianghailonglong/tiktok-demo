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
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

// CommentActionResponse 评论操作的相应结构体
type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// CommentResponse 评论信息的响应结构体
type CommentResponse struct {
	ID         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	User       User   `json:"user,omitempty"`
}

// FavoriteVideo 喜欢视频的响应结构体
type FavoriteVideo struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type FavoriteListResponse struct {
	Response
	VideoList []FavoriteVideo `json:"video_list,omitempty"`
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