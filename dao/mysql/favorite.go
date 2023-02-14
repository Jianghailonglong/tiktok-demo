package mysql

import (
	"gorm.io/gorm"
	"tiktok-demo/logger"
)

const (
	FAVORITED   = 1
	UNFAVORITED = 0
)

type Favorite struct {
	Id         int `gorm:"column:id"`
	UserId     int `gorm:"user_id"`
	VideoId    int `gorm:"video_id"`
	IsFavorite int `gorm:"is_favorite"`
}

// CheckFavorite 查询某用户是否点赞某视频
func CheckFavorite(uid int, vid int) bool {
	var total int64
	err := db.Table("favorites").Where("user_id = ? AND video_id = ? AND is_favorite = 1", uid, vid).Count(&total).Error
	if err == gorm.ErrRecordNotFound {
		return false // 点赞不存在
	}
	if total == 0 {
		return false
	}
	return true
}

// GetFavoriteVideoIdList 获取用户点赞视频id
func GetFavoriteVideoIdList(userId int) ([]int, error) {
	var favoriteVideoIdList []int
	if err := db.Table("favorites").Where("user_id=? AND is_favorite=?", userId, 1).Pluck("video_id", &favoriteVideoIdList).Error; err != nil {
		return nil, err
	}

	return favoriteVideoIdList, nil
}

// GetVideoDetail 获取视频详细信息
func GetVideoDetail(videoId int) (*Video, error) {
	var video = Video{}
	if err := db.Table("videos").Where("id=?", videoId).Find(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// GetFavorite 根据userId查找videoId是否有对应记录
func GetFavorite(userId int, videoId int) (*Favorite, error) {
	favorite := Favorite{}
	if err := db.
		Where("user_id = ?", userId).
		Where("video_id = ?", videoId).
		Take(&favorite).Error; nil != err {
		logger.Log.Warn(err.Error())
		return nil, err
	}

	return &favorite, nil
}

// UpdateFavorite 在原来的关系上修改点赞
func UpdateFavorite(favorite *Favorite, action int) error {
	// 更新失败，返回错误。
	if err := db.Model(Favorite{}).
		Where("id = ?", favorite.Id).
		Update("is_favorite", action).Error; nil != err {
		// 更新失败，打印错误日志。
		logger.Log.Error(err.Error())
		return err
	}
	// 更新成功。
	return nil
}

// AddFavorite 原表中没有点赞记录，新增一条记录
func AddFavorite(userId int, videoId int) error {
	favorite := Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: FAVORITED,
	}
	// 插入失败，返回err.
	if err := db.Select("UserId", "VideoId", "IsFavorite").Create(&favorite).Error; nil != err {
		logger.Log.Error(err.Error())
		return err
	}
	// 插入成功
	return nil
}

// UnFavorite 原表中没有点赞记录，新增一条记录
func UnFavorite(userId int, videoId int) error {
	favorite := Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: UNFAVORITED,
	}
	// 插入失败，返回err.
	if err := db.Select("UserId", "VideoId", "IsFavorite").Create(&favorite).Error; nil != err {
		logger.Log.Error(err.Error())
		return err
	}
	// 插入成功
	return nil
}

func GetFavoriteCount(videoId int) (favoriteCount int64, err error) {
	if res := db.Table("favorites").Where("video_id = ? AND is_favorite = ?", videoId, FAVORITED).Count(&favoriteCount); nil != res.Error {
		return 0, res.Error
	}
	return
}
