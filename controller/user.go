package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/logger"
	"tiktok-demo/service"
)

// Register 注册用户
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userRequest := &common.UserRequest{
		Username: username,
		Password: password,
	}
	err := c.ShouldBind(userRequest) 
	if err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	if user, token, err := service.Register(username, password); err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeInsertError,
				StatusMsg:  MsgFlags[CodeInsertError],
			},
		})
	} else {
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
			UserId: int32(user.Id),
			Token:  token,
		})
	}
}

// Login 登录用户
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userRequest := &common.UserRequest{
		Username: username,
		Password: password,
	}
	err := c.ShouldBind(userRequest)
	if err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	// 检查登录信息是否正确
	if user, token, err := service.Login(username, password); err != nil {
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeQueryError,
				StatusMsg:  MsgFlags[CodeQueryError],
			},
		})
	} else {
		c.JSON(http.StatusOK, common.UserResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
			UserId: int32(user.Id),
			Token:  token,
		})
	}
}

// UserInfo 显示用户信息
func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userInfoRequest := &common.UserInfoRequest{UserId: userIdStr}
	err := c.ShouldBind(userInfoRequest)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	userInfo, err := service.GetUserInfo(userId)
	if err != nil {
		c.JSON(http.StatusOK, common.UserInfoResponse{
			Response: common.Response{
				StatusCode: CodeQueryError,
				StatusMsg:  MsgFlags[CodeQueryError],
			},
		})
		return
	}
	c.JSON(http.StatusOK, common.UserInfoResponse{
		Response: common.Response{
			StatusCode: CodeSuccess,
			StatusMsg:  MsgFlags[CodeSuccess],
		},
		UserInfo: userInfo,
	})
}
