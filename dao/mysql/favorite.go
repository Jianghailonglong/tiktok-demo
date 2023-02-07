package mysql

import (
	"gorm.io/gorm"
	"tiktok-demo/common"
)

type Favorite struct {
	gorm.Model
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
	State   int64
}

// 查询某用户是否点赞某视频
func CheckFavorite(uid int, vid int) bool {
	var total int64
	err := db.Table("favorites").
		Where("user_id = ? AND video_id = ? AND state = 1", uid, vid).Count(&total).
		Error
	if err == gorm.ErrRecordNotFound {
		return false //关注不存在
	}
	if total == 0 {
		return false
	}
	return true
}

// 增加total_favorited
func AddTotalFavorited(HostId int) error {
	if err := db.Model(&common.User{}).
		Where("id=?", HostId).
		Update("total_favorited", gorm.Expr("total_favorited+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 减少total_favorited
func ReduceTotalFavorited(HostId int) error {
	if err := db.Model(&common.User{}).
		Where("id=?", HostId).
		Update("total_favorited", gorm.Expr("total_favorited-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 增加favorite_count
func AddFavoriteCount(HostId int) error {
	if err := db.Model(&common.User{}).
		Where("id=?", HostId).
		Update("favorite_count", gorm.Expr("favorite_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 减少favorite_count
func ReduceFavoriteCount(HostId int) error {
	if err := db.Model(&common.User{}).
		Where("id=?", HostId).
		Update("favorite_count", gorm.Expr("favorite_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func GetFavoriteList(userId int) ([]Favorite, error) {
	var favoriteList []Favorite
	if err := db.Table("favorites").Where("user_id=? AND state=?", userId, 1).Find(&favoriteList).Error; err != nil {
		return nil, err
	}
	return favoriteList, nil
}

func GetVideoDetail(videoId int) (*Video, error) {
	var video = Video{}
	if err := db.Table("videos").Where("id=?", videoId).Find(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func btoi(isFavorite bool) int64 {
	if isFavorite == true {
		return 1
	} else {
		return 0
	}
}

// 添加喜欢数量1
func FavoriteAdd(userId int, videoId int, isFavorited bool) error {
	favorite := Favorite{
		UserId:  int64(userId),
		VideoId: int64(videoId),
		State:   (btoi(isFavorited)), // 1点赞， 0未点赞
	}
	result := db.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite)
	if result.Error != nil {
		if err := db.Table("favorites").Create(&favorite).Error; err != nil {
			return err
		}
	} else {
		db.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).Update("state", favorite.State)
	}

	if isFavorited {
		if favorite.State == 0 {
			db.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			//userId的favorite_count增加
			if err := AddFavoriteCount(userId); err != nil {
				return err
			}
			//videoId对应的userId的total_favorite增加
			GuestId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}
			if err := AddTotalFavorited(int(GuestId)); err != nil {
				return err
			}
		}
	}
	return nil
}

//减少喜欢数量1
func FavoriteReduce(userId int, videoId int, isFavorited bool) error {
	favorite := Favorite{
		UserId:  int64(userId),
		VideoId: int64(videoId),
		State:   btoi(isFavorited),
	}
	result := db.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite)
	if result.Error != nil {
		if err := db.Table("favorites").Create(&favorite).Error; err != nil {
			return err
		}
	} else {
		db.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).Update("state", favorite.State)
	}

	if !isFavorited {
		if favorite.State == 0 {
			db.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			//userId的favorite_count减少
			if err := ReduceFavoriteCount(userId); err != nil {
				return err
			}
			//videoId对应的userId的total_favorite减少
			GuestId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}
			if err := ReduceTotalFavorited(int(GuestId)); err != nil {
				return err
			}
		}
	}
	return nil

}
