package mysql

import (
	"tiktok-demo/logger"
	"time"
)

type Comment struct {
	Id        int       `gorm:"column:id"`
	UserId    int       `gorm:"user_id"`
	VideoId   int       `gorm:"video_id"`
	CreatedAt time.Time `gorm:"created_at"`
	Content   string    `gorm:"content"`
	IsComment int       `gorm:"is_comment"`
}

func (Comment) TableName() string {
	return "comments"
}

const (
	COMMENTED   = 1
	UNCOMMENTED = 0
)

// GetComment 根据userId和videoId获取Comment
func GetComment(userId int, videoId int) (*Comment, error) {
	comment := Comment{}
	if err := db.
		Where("user_id = ?", userId).
		Where("video_id = ?", videoId).
		Take(&comment).Error; nil != err {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return &comment, nil
}

// AddComment 添加评论
func AddComment(userId int, videoId int, content string) (*Comment, error) {
	comment := Comment{
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
		CreatedAt: time.Now(),
		IsComment: COMMENTED,
	}

	// 插入失败，返回err
	if err := db.Create(&comment).Error; nil != err {
		logger.Log.Error(err.Error())
		return nil, err
	}
	// 插入成功
	return &comment, nil
}

func UpdateComment(comment *Comment, status int) error {
	// 更新失败，返回错误。
	values := map[string]interface{}{"is_comment": status, "content": comment.Content, "created_at": comment.CreatedAt}
	if err := db.Model(&comment).Select("is_comment", "content", "created_at").Updates(values).Error; nil != err {
		// 更新失败，打印错误日志。
		logger.Log.Error(err.Error())
		return err
	}
	// 更新成功。
	return nil
}

func GetCommentByCommentId(commentId int) (*Comment, error) {
	comment := Comment{}
	if err := db.
		Where("id = ?", commentId).Take(&comment).Error; nil != err {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return &comment, nil
}

// GetCommentList 根据videoId，获取当前视频下的评论
func GetCommentList(videoId int) ([]Comment, error) {
	var commentList []Comment
	if err := db.Table("comments").Where("video_id=? and is_comment=?", videoId, COMMENTED).Find(&commentList).Error; err != nil {
		return nil, err
	}
	return commentList, nil
}

// GetCommentCount 获取视频评论数
func GetCommentCount(videoId int) (commentCount int64, err error) {
	if res := db.Table("comments").Where("video_id = ? AND is_comment = ?", videoId, COMMENTED).Count(&commentCount); nil != res.Error {
		return 0, res.Error
	}
	return
}
