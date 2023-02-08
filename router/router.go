package router

import (
	"github.com/gin-gonic/gin"
	"tiktok-demo/controller"
	"tiktok-demo/middleware/jwt"
)

func InitRouters(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200,"连接成功!")
	})

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/user/", jwt.AuthInHeader(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// 视频流接口不限制登录状态，登录和非登录状态对视频流内容获取有不同处理
	// apiRouter.GET("/feed/", jwt.AuthWithoutLimitLoginStatus(), controller.Feed)
	// 投稿视频接口token在body里
	// apiRouter.POST("/publish/action/", jwt.AuthInBody(), controller.Publish)
	// apiRouter.GET("/publish/list/", jwt.AuthInHeader(), controller.PublishList)

	// extra apis - I
	//apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	//apiRouter.GET("/favorite/list/", controller.FavoriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", jwt.AuthInHeader(), controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	//apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/",jwt.AuthInHeader() ,controller.MessageChat)
	apiRouter.POST("/message/action/",jwt.AuthInHeader() ,controller.MessageAction)

	
	//聊天模块(方案一:websocket)
	//websocket不能够使用gin的中间件 T_T
	//所以要自己实现鉴权
	r.GET("/douyin/chat/ws",controller.WsChatHandler)
}
