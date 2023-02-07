package mysql

import (
	"gorm.io/gorm"
	"tiktok-demo/common"
	"time"
)

type Video struct {
	gorm.Model
	AuthorId      int64  `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	Title         string `json:"title"`
}

type Comment struct {
	gorm.Model
	VideoId int64  `json:"video_id,omitempty"`
	UserId  int64  `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}

func GetCommentList(videoId int) ([]Comment, error) {
	var commentList []Comment
	if err := db.Table("comments").Where("video_id=?", videoId).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

func PostComment(comment Comment) error {
	if err := db.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentId int) error {
	if err := db.Table("comments").Where("id = ?", commentId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func AddCommentCount(videoId int) error {
	if err := db.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

func ReduceCommentCount(videoId int) error {
	if err := db.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

func GetVideoAuthor(videoId int) (int64, error) {
	var video Video
	if err := db.Table("videos").Where("id = ?", videoId).Find(&video).Error; err != nil {
		return int64(video.ID), err
	}
	return int64(video.ID), nil
}

func GetUser(userId int) (common.User, error) {
	var user common.User
	if err := db.Model(&common.User{}).Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
