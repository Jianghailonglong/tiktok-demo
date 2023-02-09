package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/middleware/jwt"
	"tiktok-demo/service"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	userIDRaw, _ := c.Get(jwt.ContextUserIDKey)
	userId, ok := userIDRaw.(int)

	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, err1 := strconv.ParseInt(videoIdStr, 10, 10)

	if (actionType != "1" && actionType != "2") || err1 != nil || !ok {
		c.JSON(http.StatusOK, common.CommentActionResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
			Comment: nil,
		})
		return
	}

	if actionType == "1" {
		text := c.Query("comment_text")
		comment := service.CommentPost(userId, int(videoId), text)
		c.JSON(http.StatusOK, common.CommentActionResponse{
			Response: common.Response{
				StatusCode: CodeSuccess,
				StatusMsg:  MsgFlags[CodeSuccess],
			},
			Comment: comment},
		)
		return
	} else if actionType == "2" {
		commentIdStr := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, common.CommentActionResponse{
				Response: common.Response{
					StatusCode: CodeInvalidParams,
					StatusMsg:  MsgFlags[CodeInvalidParams],
				},
				Comment: nil,
			})
		}
		if err = service.CommentDelete(int(commentId)); nil == err {
			c.JSON(http.StatusOK, common.CommentActionResponse{
				Response: common.Response{
					StatusCode: CodeSuccess,
					StatusMsg:  MsgFlags[CodeSuccess],
				},
				Comment: nil},
			)
			return
		}
	}
}

// CommentList 获取评论列表
func CommentList(c *gin.Context) {
	userIDRaw, _ := c.Get(jwt.ContextUserIDKey)
	userId, ok := userIDRaw.(int)
	videoIdStr := c.Query("video_id")
	videoId, err1 := strconv.ParseInt(videoIdStr, 10, 10)
	if err1 != nil || !ok {
		c.JSON(http.StatusOK, common.CommentListResponse{
			Response: common.Response{
				StatusCode: CodeInvalidParams,
				StatusMsg:  MsgFlags[CodeInvalidParams],
			},
			CommentList: nil,
		})
		return
	}

	commentListResponse := service.GetCommentList(userId, int(videoId))

	if commentListResponse != nil {
		commentListResponse.Response = common.Response{
			StatusCode: CodeSuccess,
			StatusMsg:  MsgFlags[CodeSuccess],
		}
		c.JSON(http.StatusOK, commentListResponse)
	} else {
		commentListResponse.Response = common.Response{
			StatusCode: CodeQueryError,
			StatusMsg:  MsgFlags[CodeQueryError],
		}
		c.JSON(http.StatusOK, commentListResponse)
	}
}
