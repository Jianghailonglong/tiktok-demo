package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/service"
)

// 点赞视频方法
func Favorite(c *gin.Context) {
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}

	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseInt(actionTypeStr, 10, 10)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 10)

	err := service.FavoriteAction(userId, uint(videoId), uint(actionType))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "操作成功！",
		})
	}
}

// 获取列表方法
func FavoriteList(c *gin.Context) {
	getUserId, _ := c.Get("user_id")
	var userIdHost int
	if v, ok := getUserId.(int); ok {
		userIdHost = v
	}
	userIdStr := c.Query("user_id")
	userId, _ := strconv.ParseUint(userIdStr, 10, 10)
	userIdNew := int(userId)
	if userIdNew == 0 {
		userIdNew = userIdHost
	}

	videoList, err := service.FavoriteList(userIdNew)
	videoListNew := make([]common.FavoriteVideo, 0)
	for _, m := range videoList {
		var author = common.User{}
		var getAuthor = common.User{}
		getAuthor, err := service.GetUser(int(m.AuthorId))
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "找不到作者！",
			})
			c.Abort()
			return
		}
		isfollowing := service.IsFollowing(userIdHost, uint(m.AuthorId)) //参数类型、错误处理
		isfavorite := service.CheckFavorite(userIdHost, int(m.ID))

		author.Id = getAuthor.Id
		author.Name = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing
		var video = common.FavoriteVideo{}
		video.Id = int64(m.ID)
		video.Author = author
		video.PlayUrl = m.PlayUrl
		video.CoverUrl = m.CoverUrl
		video.FavoriteCount = m.FavoriteCount
		video.CommentCount = m.CommentCount
		video.IsFavorite = isfavorite
		video.Title = m.Title

		videoListNew = append(videoListNew, video)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, common.FavoriteListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, common.FavoriteListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			VideoList: videoListNew,
		})
	}
}
