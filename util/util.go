package util

import (
	"sort"
	"tiktok-demo/common"
)

// SortVideoStrategy 策略模式实现按点赞数和评论数排序视频
type SortVideoStrategy interface {
	sortVideoList([]common.Video)
}

type SortVideoByFavoriteCount struct {
}

func (f *SortVideoByFavoriteCount) sortVideoList(videoList []common.Video) {
	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].FavoriteCount > videoList[j].FavoriteCount
	})
}

type SortVideoByCommentCount struct {
}

func (f *SortVideoByCommentCount) sortVideoList(videoList []common.Video) {
	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].CommentCount > videoList[j].CommentCount
	})
}

// SortVideoContext 上下文类
type SortVideoContext struct {
	strategy SortVideoStrategy
}

func (f *SortVideoContext) SetSortVideoStrategy(s SortVideoStrategy) {
	f.strategy = s
}

func (f *SortVideoContext) SortVideo(videoList []common.Video) {
	f.strategy.sortVideoList(videoList)
}
