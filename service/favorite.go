package service

import (
	"tiktok-demo/dao/mysql"
)

// 查询某用户是否点赞某视频
func CheckFavorite(uid int, vid int) bool {
	return mysql.CheckFavorite(uid, vid)
}

//点赞操作
func FavoriteAction(userId uint, videoId uint, actionType uint) (err error) {
	if actionType == 1 {
		return mysql.FavoriteAdd(int(userId), int(videoId), true)
	}
	return mysql.FavoriteReduce(int(userId), int(videoId), false)
}

// 获取点赞列表
func FavoriteList(userId int) ([]mysql.Video, error) {
	favoriteList, err := mysql.GetFavoriteList(userId)
	if err != nil {
		return nil, err
	}
	videoList := make([]mysql.Video, 0)
	for _, m := range favoriteList {
		video, err := mysql.GetVideoDetail(int(m.VideoId))
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, *video)
	}
	return videoList, nil
}
