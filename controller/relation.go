package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/service"
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
		if err := service.SubscribeUser(userId, toUserId); nil == err {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 0,
				StatusMsg:  "关注成功",
			})
			return
		}
	} else if actionType == 2 {
		if err := service.UnsubscribeUser(userId, toUserId); nil == err {
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

// /follow/list 获取关注列表
func FollowList(c *gin.Context) {
	userId := c.GetInt(ContextUserIDKey)

	if followList, err := service.GetFollowList(int64(userId)); nil == err {
		c.JSON(http.StatusOK, common.FollowListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "获取关注列表成功",
			},
			UserList: followList,
		})
		return
	}

	c.JSON(http.StatusOK, common.FollowListResponse{
		Response: common.Response{
			StatusCode: -1,
			StatusMsg:  "获取关注列表失败",
		},
		UserList: nil,
	})

}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	userId := c.GetInt(ContextUserIDKey)

	if followerList, err := service.GetFollowerList(int64(userId)); nil == err {
		c.JSON(http.StatusOK, common.FollowerListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "获取粉丝列表成功",
			},
			UserList: followerList,
		})
		return
	}

	c.JSON(http.StatusOK, common.FollowerListResponse{
		Response: common.Response{
			StatusCode: -1,
			StatusMsg:  "获取粉丝列表失败",
		},
		UserList: nil,
	})
}
