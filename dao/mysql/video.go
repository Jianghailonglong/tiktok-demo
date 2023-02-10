package mysql

import (
	"errors"
	"strconv"
	"tiktok-demo/conf"
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
func InsertVideo(userID int, title, videoName, imageName string) (err error) {
	video := Video{
		AuthorId:    userID,
		PlayUrl:     conf.Config.MinioConfig.Video.URL + videoName, // 访问nginx接口进行代理，易于后续实现分布式minio扩展
		CoverUrl:    conf.Config.MinioConfig.Image.URL + imageName,
		Title:       title,
		PublishTime: time.Now().Unix(),
	}
	res := db.Create(&video)
	if res.Error != nil {
		return errors.New("InsertVideo插入失败")
	}
	return
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
