package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/middleware/jwt"
	"tiktok-demo/service"
)

// FavoriteAction 点赞视频方法
func FavoriteAction(c *gin.Context) {
	getUserId, _ := c.Get(jwt.ContextUserIDKey)
	var userId int
	if v, ok := getUserId.(int); ok {
		userId = v
	}

	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseInt(actionTypeStr, 10, 10)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 10)

	err := service.FavoriteAction(userId, int(videoId), int(actionType))
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

// FavoriteList 获取喜欢视频列表方法
func FavoriteList(c *gin.Context) {
	userIDRaw := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.FavoriteListResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}

	favoriteListResponse, err := service.FavoriteList(int(userID))

	if err != nil {
		c.JSON(http.StatusBadRequest, common.FavoriteListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			VideoList: nil,
		})
	} else {
		favoriteListResponse.Response = common.Response{
			StatusCode: 0,
			StatusMsg:  "已找到列表！",
		}
		c.JSON(http.StatusOK, favoriteListResponse)
	}
}
