package service

import (
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
)

// 获取指定videoId的评论表
func GetCommentList(videoId int) ([]mysql.Comment, error) {
	return mysql.GetCommentList(videoId)
}

// 发布评论
func PostComment(comment mysql.Comment) error {
	return mysql.PostComment(comment)
}

// 删除指定commentId的评论
func DeleteComment(commentId int) error {
	return mysql.DeleteComment(commentId)
}

// 增加评论数量1
func AddCommentCount(videoId int) error {

	return mysql.AddCommentCount(videoId)
}

// 减少评论数量1
func ReduceCommentCount(videoId int) error {
	return mysql.ReduceCommentCount(videoId)
}

// 根据用户id获取用户信息
func GetUser(userId int) (common.User, error) {
	return mysql.GetUser(userId)
}

// 判断HostId是否关注GuestId
//func IsFollowing(HostId uint, GuestId uint) bool {
//	var relationExist = &common.Following{}
//	err := mysql.Sql.Model(&common.Following{}).
//		Where("host_id=? AND guest_id=?", HostId, GuestId).
//		First(&relationExist).Error
//	if err == gorm.ErrRecordNotFound {
//		return false
//	}
//	return true
//}
