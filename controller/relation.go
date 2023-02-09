package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/middleware/jwt"
	"tiktok-demo/service"
)

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	userId := c.GetInt(jwt.ContextUserIDKey)
	toUserId, err1 := strconv.Atoi(c.Query("to_user_id"))
	actionType, err2 := strconv.Atoi(c.Query("action_type"))
	if (err1 != nil || err2 != nil) || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: CodeInvalidParams,
			StatusMsg:  MsgFlags[CodeInvalidParams],
		})
		return
	}

	if actionType == 1 {
		if err := service.SubscribeUser(userId, toUserId); nil == err {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			})
			return
		}
	} else if actionType == 2 {
		if err := service.UnsubscribeUser(userId, toUserId); nil == err {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			})
			return
		}
	}

	c.JSON(http.StatusOK, common.Response{
		StatusCode: CodeServerBusy,
		StatusMsg:  MsgFlags[CodeServerBusy],
	})
}

// FollowList /follow/list 获取关注列表
func FollowList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.FollowListResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}

	if followList, err := service.GetFollowList(userId); nil == err {
		c.JSON(http.StatusOK, common.FollowListResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
			UserList: followList,
		})
		return
	}

	c.JSON(http.StatusOK, common.FollowListResponse{
		Response: common.Response{
			StatusCode: CodeServerBusy,
			StatusMsg:  MsgFlags[CodeServerBusy],
		},
		UserList: nil,
	})

}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.FollowerListResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}

	if followerList, err := service.GetFollowerList(userId); nil == err {
		c.JSON(http.StatusOK, common.FollowerListResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
			UserList: followerList,
		})
		return
	}

	c.JSON(http.StatusOK, common.FollowerListResponse{
		Response: common.Response{
			StatusCode: CodeServerBusy,
			StatusMsg:  MsgFlags[CodeServerBusy],
		},
		UserList: nil,
	})
}
