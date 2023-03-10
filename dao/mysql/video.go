package mysql

import (
	"errors"
	"strconv"
	"tiktok-demo/logger"
	"time"
)

type Video struct {
	Id          int    `gorm:"column:id"`
	AuthorId    int    `gorm:"column:author_id"`
	PlayUrl     string `gorm:"column:play_url"`
	CoverUrl    string `gorm:"column:cover_url"`
	Title       string `gorm:"column:title"`
	PublishTime int64  `gorm:"column:publish_time"`
}

func (Video) TableName() string {
	return "videos"
}

// GetFeed 获取视频流
func GetFeed(latestTimeRaw string) (videos []Video, err error) {
	latestTime, err := strconv.ParseInt(latestTimeRaw, 10, 64)
	if err != nil {
		logger.Log.Error("strconv.ParseInt failed")
		return
	}
	res := db.Table("videos").Order("publish_time desc").
		Where("publish_time <= ?", latestTime).Limit(5).Scan(&videos)
	if res.Error != nil {
		logger.Log.Error("GetFeed获取视频流失败")
		return videos, res.Error
	}
	return
}

// InsertVideo 插入视频
func InsertVideo(userID int, title, playUrl, coverUrl string) (videoId int, err error) {
	video := Video{
		AuthorId:    userID,
		PlayUrl:     playUrl, // 访问nginx接口进行代理，易于后续实现分布式minio扩展
		CoverUrl:    coverUrl,
		Title:       title,
		PublishTime: time.Now().Unix(),
	}
	res := db.Create(&video)
	if res.Error != nil {
		return 0, errors.New("InsertVideo插入失败")
	}
	return video.Id, nil
}

// DeleteVideo 删除视频
func DeleteVideo(videoId int) (err error) {
	res := db.Where("id = ?", videoId).Delete(&Video{})
	if res.Error != nil {
		return errors.New("DeleteVideo删除失败")
	}
	return nil
}

// GetVideoListByUserID 根据userID获取全部视频列表
func GetVideoListByUserID(userID int64) (videos []Video, err error) {
	res := db.Where("author_id = ?", userID).Find(&videos)
	if res.Error != nil {
		logger.Log.Error("GetVideoListByUserID获取视频失败")
		return videos, res.Error
	}
	return
}
