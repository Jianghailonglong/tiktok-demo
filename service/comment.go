package service

import (
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
)

// CommentPost 发布评论
func CommentPost(userId int, videoId int, text string) *common.Comment {
	// 1、需要创建评论
	comment, _ := mysql.AddComment(userId, videoId, text)
	if comment == nil {
		logger.Log.Error("mysql.AddComment failed")
		return nil
	}
	// 2、获取用户信息
	if user, err := GetCommonUserInfoById(int64(userId), int64(userId)); nil == err {
		return &common.Comment{
			Id:         int64(comment.Id),
			User:       user,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("2006-01-02"),
		}
	} else {
		logger.Log.Error("获取用户信息失败")
		return nil
	}
}

// CommentDelete 删除评论
func CommentDelete(commentId int) error {
	return mysql.DeleteComment(commentId)
}

// GetCommentList 获取指定videoId的评论表
func GetCommentList(userId, videoId int) *common.CommentListResponse {
	// 1、从数据库中获取视频下的评论表
	commentList, _ := mysql.GetCommentList(videoId)
	if commentList == nil {
		return nil
	}
	// 2、获取评论用户信息
	n := len(commentList)
	userList := make([]common.User, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(int64(userId), int64(commentList[i].UserId)); nil == err {
				userList[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}
	wg.Wait()
	return assembleCommentList(commentList, userList)
}

func assembleCommentList(comments []mysql.Comment, users []common.User) (commentListResponse *common.CommentListResponse) {
	commentListResponse = new(common.CommentListResponse)
	commentListResponse.CommentList = make([]common.Comment, len(comments))
	for i := 0; i < len(comments); i++ {
		commentListResponse.CommentList[i].Id = int64(comments[i].Id)
		commentListResponse.CommentList[i].User = users[i]
		commentListResponse.CommentList[i].Content = comments[i].Content
		commentListResponse.CommentList[i].CreateDate = comments[i].CreatedAt.Format("2006-01-02")
	}
	return
}
