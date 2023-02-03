package controller

import (
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/service"
	"github.com/gin-gonic/gin"
)

const (
	// ContextUserIDKey 验证token成功后，将userID加入到context上下文中
	ContextUserIDKey = "userID"
)

// Action 关注操作
func RelationAction(c *gin.Context) {
	userId := c.GetInt(ContextUserIDKey)
	toUserId, err1 := strconv.Atoi(c.Query("to_user_id"))
	actionType, err2 := strconv.Atoi(c.Query("action_type"))
	if (err1 != nil || err2 != nil) || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  "关注参数错误",
		})
		return
	}

	if actionType == 1 {
		if flag, _ := service.SubscribeUser(userId, toUserId); flag {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 0,
				StatusMsg:  "关注成功",
			})
			return
		}
	} else if actionType == 2 {
		if flag, _ := service.UnsubscribeUser(userId, toUserId); flag {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 0,
				StatusMsg:  "取消关注成功",
			})
			return
		}
	}

	c.JSON(http.StatusOK, common.Response{
		StatusCode: -1,
		StatusMsg:  "服务器内部错误",
	})
}
