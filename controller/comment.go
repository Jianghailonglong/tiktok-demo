package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/service"
)

// 评论操作
func CommentAction(c *gin.Context) {
	rawUserId, _ := c.Get("user_id")

	var userId int
	if id, ok := rawUserId.(int); ok {
		userId = id
	}

	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 10)
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "Unsupported actionType",
		})
		c.Abort()
		return
	}

	switch actionType {
	case "1":
		text := c.Query("comment_text")
		CommentPost(c, userId, text, videoId)
	case "2":
		commentStrId := c.Query("comment_id")
		commentId, _ := strconv.ParseUint(commentStrId, 10, 10)
		CommentDelete(c, int(videoId), int(commentId))
	default:
		return
	}
}

// 发布评论
func CommentPost(c *gin.Context, userId int, text string, videoId int64) {
	newComment := mysql.Comment{
		VideoId: videoId,
		UserId:  int64(userId),
		Content: text,
	}

	err1 := service.PostComment(newComment)

	err2 := service.AddCommentCount(int(videoId))

	getUser, err3 := service.GetUser(userId)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to post comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, common.CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "post the comment successfully",
		},
		Comment: common.CommentResponse{
			ID:         int64(newComment.ID),
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: common.User{
				Id:            getUser.Id,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				//IsFollow:      service.IsFollowing(userId, videoAuthor),
			},
		},
	})
}

// 删除评论
func CommentDelete(c *gin.Context, videoId int, commentId int) {

	err1 := service.DeleteComment(commentId);

	err2 := service.ReduceCommentCount(videoId);

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to delete comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "delete the comment successfully",
	})
}

// 获取评论列表
func CommentList(c *gin.Context) {
	rawUserId, _ := c.Get("user_id")

	var userId uint
	if id, ok := rawUserId.(uint); ok {
		userId = id
	}

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 10)

	commentList, err := service.GetCommentList(int(videoId))

	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}
	var responseCommentList []common.CommentResponse
	for i := 0; i < len(commentList); i++ {
		getUser, err1 := service.GetUser(int(commentList[i].UserId))

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseComment := common.CommentResponse{
			ID:         int64(commentList[i].ID),
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: common.User{
				Id:            getUser.Id,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, commentList[i].ID),
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	//响应返回
	c.JSON(http.StatusOK, common.CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the comment list.",
		},
		CommentList: responseCommentList,
	})

}
