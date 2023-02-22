package controller

import (
	"net/http"
	"tiktok-demo/common"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/jwt"
	"tiktok-demo/service"

	"github.com/gin-gonic/gin"
)

func MessageAction(c *gin.Context) {
	// token := c.Query("token")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	content := c.Query("content")
	chatMessageRequest := &common.ChatMessageRequest{
		// Token:      token,
		ToUserId:   toUserId,
		ActionType: actionType,
		Content:    content,
	}
	err := c.ShouldBind(chatMessageRequest)
	if err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.ChatMessageResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	//聊天消息数据插入
	userIdRaw, _ := c.Get(jwt.ContextUserIDKey)
	userId, _ := userIdRaw.(int)
	if err = service.AddChatRecord(userId, toUserId, actionType, content); err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.ChatMessageResponse{
			Response: common.Response{
				StatusCode: CodeInsertError,
				StatusMsg:  MsgFlags[CodeInsertError],
			},
		})
	} else {
		c.JSON(http.StatusOK, common.ChatMessageResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
		})
	}
}

func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	chatHistoryRequest := &common.ChatHistoryRequest{
		Token:    token,
		ToUserId: toUserId,
	}
	//用于验证参数合法性
	err := c.ShouldBind(chatHistoryRequest)
	if err != nil {
		logger.Log.Error(err.Error())
		c.JSON(http.StatusOK, common.ChatHistoryResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
		})
		return
	}
	//获取消息列表
	userIdRaw, _ := c.Get(jwt.ContextUserIDKey)
	userId, _ := userIdRaw.(int)
	//messageList,err:=service.GetChatRecordList(userId,toUserId)
	messageList, err := service.GetChatUnreadList(userId, toUserId)
	if err != nil {
		c.JSON(http.StatusOK, common.ChatHistoryResponse{
			Response: common.Response{
				StatusCode: CodeQueryError,
				StatusMsg:  MsgFlags[CodeQueryError],
			},
		})
		return
	}
	c.JSON(http.StatusOK, common.ChatHistoryResponse{
		Response: common.Response{
			StatusCode: CodeSuccess,
			StatusMsg:  MsgFlags[CodeSuccess],
		},
		MessageList: messageList,
	})
}
