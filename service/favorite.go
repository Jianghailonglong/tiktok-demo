package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/kafka"
)

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context, userId int, videoId int, actionType int) (err error) {
	if actionType == 1 {
		return FavoriteVideo(c, userId, videoId)
	}
	return UnFavoriteVideo(c, userId, videoId)
}

// FavoriteVideo 点赞视频
func FavoriteVideo(c *gin.Context, userId int, videoId int) (err error) {
	// 1、更新mysql
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation == nil {
		err = mysql.AddFavorite(userId, videoId)
	} else {
		err = mysql.UpdateFavorite(relation, mysql.FAVORITED)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2、将更新信息传给kafka，由kafka来重试删除或者更新redis
	// kafka点赞消息格式 key: videoId
	// value: userId:del，删除，value: userId:add，添加， value: cnt，添加视频点赞数
	err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoId), strconv.Itoa(userId)+":del")
	if err != nil {
		logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
		return err
	}
	return nil
}

// UnFavoriteVideo 取消点赞
func UnFavoriteVideo(c *gin.Context, userId int, videoId int) (err error) {
	// 1、更新mysql
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation != nil {
		err = mysql.UpdateFavorite(relation, mysql.UNFAVORITED)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2、将更新信息传给kafka，由kafka来重试删除或者更新redis
	// value: userId:del，删除 value: userId:add，添加
	err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoId), strconv.Itoa(userId)+":del")
	if err != nil {
		logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
		return err
	}
	return nil
}

// FavoriteList 获取点赞视频列表
func FavoriteList(c *gin.Context, userId int) (favoriteListResponse *common.FavoriteListResponse, err error) {
	// 1、获取视频列表id
	// 1) redis获取
	videoIdList, err := redis.GetFavoriteVideoList(c, userId)
	if err != nil || len(videoIdList) == 0 {
		// 2) 从mysql获取视频列表
		videoIdList, err = mysql.GetFavoriteVideoIdList(userId)
		if err != nil {
			return nil, err
		}
		// 3）将mysql获取到的消息传给kafka去更新redis
		for _, videoId := range videoIdList {
			err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoId), strconv.Itoa(userId)+":add")
			if err != nil {
				logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
				return nil, err
			}
		}
	}

	videoList := make([]mysql.Video, len(videoIdList))
	size := len(videoList)
	var wg sync.WaitGroup
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if video, err := mysql.GetVideoDetail(videoIdList[i]); err == nil {
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
	// 1）从redis获取
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			// 计数redis可能有偏差，但影响不大
			if cnt, err := redis.GetVideoFavoriteCount(c, videoIdList[i]); nil == err && cnt != 0 {
				videoLikeCntsList[i] = cnt
			} else {
				logger.Log.Error("从redis获取视频点赞次数失败")
				// 2）从mysql获取
				if cnt, err := mysql.GetFavoriteCount(videoIdList[i]); nil == err {
					videoLikeCntsList[i] = int(cnt)
					// 3）传给mysql更新点赞数
					err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoIdList[i]), strconv.Itoa(int(cnt)))
					if err != nil {
						logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
					}
				} else {
					logger.Log.Error("从mysql获取点赞次数失败")
				}
			}
		}(i)
	}
	// 4、获取视频评论数
	videoCommentCntsList := make([]int, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := mysql.GetCommentCount(videoIdList[i]); nil == err {
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
		favoriteListResponse.VideoList[i].CommentCount = int64(videoCommentCntsList[i])
		favoriteListResponse.VideoList[i].IsFavorite = isFavoriteVideoList[i]
	}
	return
}
