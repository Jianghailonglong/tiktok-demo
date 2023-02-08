package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/jwt"
	"tiktok-demo/service"
	"time"
)

// Feed 视频流
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	if latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	}
	userIDRaw, _ := c.Get(jwt.ContextUserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusOK, common.FeedResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	feedResponse, err := service.GetFeed(latestTime, userID)
	if err != nil {
		logger.Log.Error("service.GetFeed failed", zap.Any("error", err))
		c.JSON(http.StatusOK, common.FeedResponse{
			Response: common.Response{
				StatusCode: CodeQueryError,
				StatusMsg:  MsgFlags[CodeQueryError],
			},
		})
		return
	}
	// 拼凑信息返回响应内容
	feedResponse.Response = common.Response{
		StatusCode: CodeSuccess,
		StatusMsg:  MsgFlags[CodeSuccess],
	}
	c.JSON(http.StatusOK, feedResponse)
}

// Publish 投稿视频
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	file, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	userIDRaw, _ := c.Get(jwt.ContextUserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	var inputBuffer bytes.Buffer
	_, err = io.Copy(&inputBuffer, file)
	if err != nil {
		logger.Log.Error("io.Copy failed", zap.Any("error", err.Error()))
		c.JSON(http.StatusOK, common.Response{
			StatusCode: CodeInsertError,
			StatusMsg:  MsgFlags[CodeInsertError],
		})
		return
	}
	err = service.PublishVideo(userID, title, &inputBuffer)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: CodeInsertError,
			StatusMsg:  MsgFlags[CodeInsertError],
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: CodeSuccess,
			StatusMsg:  MsgFlags[CodeSuccess],
		})
	}
}

// PublishList 用户的视频发布列表
func PublishList(c *gin.Context) {
	userIDRaw := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	videoPublishListResponse, err := service.PublishList(userID)
	if err != nil {
		logger.Log.Error("service.PublishList failed", zap.Any("error", err))
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeQueryError,
				StatusMsg:  MsgFlags[CodeQueryError],
			},
		})
		return
	}
	// 拼凑信息返回响应内容
	videoPublishListResponse.Response = common.Response{
		StatusCode: CodeSuccess,
		StatusMsg:  MsgFlags[CodeSuccess],
	}
	c.JSON(http.StatusOK, videoPublishListResponse)
}
