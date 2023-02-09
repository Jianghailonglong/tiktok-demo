package service

import (
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
)

// FavoriteAction 点赞操作
func FavoriteAction(userId int, videoId int, actionType int) (err error) {
	if actionType == 1 {
		return FavoriteVideo(userId, videoId)
	}
	return UnFavoriteVideo(userId, videoId)
}

// FavoriteVideo 点赞视频
func FavoriteVideo(userId int, videoId int) error {
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation == nil {
		return mysql.AddFavorite(userId, videoId)
	} else {
		return mysql.UpdateFavorite(relation, mysql.FAVORITED)
	}
}

// UnFavoriteVideo 取消点赞
func UnFavoriteVideo(userId int, videoId int) error {
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation == nil {
		return nil
	} else {
		return mysql.UpdateFavorite(relation, mysql.UNFAVORITED)
	}
}

// FavoriteList 获取点赞视频列表
func FavoriteList(userId int) (favoriteListResponse *common.FavoriteListResponse, err error) {
	// 1、获取视频列表
	favoriteList, err := mysql.GetFavoriteList(userId)
	if err != nil {
		return nil, err
	}
	videoList := make([]mysql.Video, len(favoriteList))
	size := len(videoList)
	var wg sync.WaitGroup
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if video, err := mysql.GetVideoDetail(favoriteList[i].VideoId); err == nil {
				videoList[i] = *video
			} else {
				logger.Log.Error("获取视频信息失败")
			}
		}(i)
	}
	wg.Wait()
	wg.Add(size * 3)
	// 2、获取用户信息
	resUsers := make([]common.User, size)

	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(int64(userId), int64(videoList[i].AuthorId)); nil == err {
				resUsers[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}
	// 3、获取视频点赞数
	videoLikeCntsList := make([]int, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := mysql.GetFavoriteCount(favoriteList[i].VideoId); nil == err {
				videoLikeCntsList[i] = int(cnt)
			} else {
				logger.Log.Error("获取点赞次数失败")
			}
		}(i)
	}
	// 4、获取视频评论数
	videoCommentCntsList := make([]int, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := mysql.GetCommentCount(favoriteList[i].VideoId); nil == err {
				videoCommentCntsList[i] = int(cnt)
			} else {
				logger.Log.Error("获取评论次数失败")
			}
		}(i)
	}
	// 5、用户的点赞视频，点赞关系
	isFavoriteVideoList := make([]bool, size)
	for i := 0; i < size; i++ {
		isFavoriteVideoList[i] = true
	}
	wg.Wait()
	favoriteListResponse = assembleFavoriteList(videoList, resUsers, videoLikeCntsList, videoCommentCntsList,
		isFavoriteVideoList)
	return
}

// assembleFavoriteList 组装视频流信息
func assembleFavoriteList(videos []mysql.Video, resUsers []common.User, videoLikeCntsList []int,
	videoCommentCntsList []int, isFavoriteVideoList []bool) (favoriteListResponse *common.FavoriteListResponse) {
	favoriteListResponse = new(common.FavoriteListResponse)
	favoriteListResponse.VideoList = make([]common.Video, len(videos))
	for i := 0; i < len(videos); i++ {
		favoriteListResponse.VideoList[i].Id = int64(videos[i].Id)
		favoriteListResponse.VideoList[i].PlayUrl = videos[i].PlayUrl
		favoriteListResponse.VideoList[i].CoverUrl = videos[i].CoverUrl
		favoriteListResponse.VideoList[i].Title = videos[i].Title
		favoriteListResponse.VideoList[i].Author = resUsers[i]
		favoriteListResponse.VideoList[i].FavoriteCount = int64(videoLikeCntsList[i])
		favoriteListResponse.VideoList[i].FavoriteCount = int64(videoCommentCntsList[i])
		favoriteListResponse.VideoList[i].IsFavorite = isFavoriteVideoList[i]
	}
	return
}
