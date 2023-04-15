package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/batcher"
	"tiktok-demo/middleware/kafka"
	"time"
)

const (
	batcherSize     = 100
	batcherBuffer   = 100
	batcherWorker   = 10
	batcherInterval = time.Second
)

type FavoriteLogic struct {
	batcher *batcher.Batcher
}

func init() {
	c := cron.New()
	// 每天3点定时刷新0 3 * * *
	_, err := c.AddFunc("0 3 * * *", func() {
		allFavorites, err := mysql.GetAllFavoriteRelationList()
		if err != nil {
			logger.Log.Error("mysql.GetAllFavoriteRelationList failed", zap.Any("error", err))
			return
		}
		videoIdCnt := make(map[int]int, 0)
		ctx := context.Background()
		for _, fav := range allFavorites {
			videoIdCnt[fav.VideoId]++
			err = redis.AddFavorite(ctx, fav.UserId, []interface{}{fav.VideoId})
			if err != nil {
				logger.Log.Error("redis.AddFavorite failed", zap.Any("error", err))
				return
			}
		}
		for videoId, cnt := range videoIdCnt {
			err := redis.SetVideoFavoriteCount(ctx, videoId, cnt)
			if err != nil {
				logger.Log.Error("redis.SetVideoFavoriteCount failed", zap.Any("error", err))
				return
			}
		}
	})
	if err != nil {
		logger.Log.Error("cron.AddFunc failed", zap.Any("error", err))
		return
	}
	c.Start()
}

func NewFavoriteLogic() *FavoriteLogic {
	f := &FavoriteLogic{}
	// 利用batcher进行批处理
	b := batcher.New(
		batcher.WithSize(batcherSize),
		batcher.WithBuffer(batcherBuffer),
		batcher.WithWorker(batcherWorker),
		batcher.WithInterval(batcherInterval),
	)
	// 需要定义Sharding和Do方法
	b.Sharding = func(key string) int {
		pid, _ := strconv.ParseInt(key, 10, 64)
		return int(pid) % batcherWorker
	}
	b.Do = func(ctx context.Context, vals map[string][]string) {
		if err = kafka.FavoriteClient.SendMessages(vals); err != nil {
			logger.Log.Error("kafka.FavoriteClient.SendMessages failed", zap.Any("error", err))
		}
	}
	f.batcher = b
	f.batcher.Start()
	return f
}

// FavoriteAction 点赞操作
func (f *FavoriteLogic) FavoriteAction(userId int, videoId int, actionType int) (err error) {
	if actionType == 1 {
		return f.FavoriteVideo(userId, videoId)
	}
	return f.UnFavoriteVideo(userId, videoId)
}

// FavoriteVideo 点赞视频
func (f *FavoriteLogic) FavoriteVideo(userId int, videoId int) (err error) {
	// 1、将消息传给kafka
	// kafka点赞消息格式 key: videoId
	// value: userId:del，删除，value: userId:add，添加， value: cnt，添加视频点赞数
	err = f.batcher.Add(strconv.Itoa(videoId), strconv.Itoa(userId)+":add")
	if err != nil {
		logger.Log.Error("f.batcher.Add failed", zap.Any("error", err))
		return err
	}
	return nil
}

// UnFavoriteVideo 取消点赞
func (f *FavoriteLogic) UnFavoriteVideo(userId int, videoId int) (err error) {
	// 1、将信息传给kafka
	// value: userId:del，删除 value: userId:add，添加
	err = f.batcher.Add(strconv.Itoa(videoId), strconv.Itoa(userId)+":del")
	if err != nil {
		logger.Log.Error("f.batcher.Add failed", zap.Any("error", err))
		return err
	}
	return nil
}

// FavoriteList 获取点赞视频列表
func FavoriteList(c *gin.Context, userId int) (favoriteListResponse *common.FavoriteListResponse, err error) {
	// 1、获取视频列表id
	// 1) redis获取用户喜欢视频id列表
	videoIdList, err := redis.GetFavoriteVideoList(c, userId)
	if err != nil {
		logger.Log.Error("redis.GetFavoriteVideoList failed")
		// return nil, err
	} else {
		videoIdList, err = mysql.GetFavoriteVideoIdList(userId)
		if err != nil {
			return nil, err
		}
	}
	// 2）获取video详细信息
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
	// 从redis获取
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := redis.GetVideoFavoriteCount(c, videoIdList[i]); nil == err {
				videoLikeCntsList[i] = cnt
			} else {
				logger.Log.Error("从redis获取视频点赞次数失败")
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
