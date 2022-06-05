package main

import (
	"douyin-server/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("")

	// basic apis
	apiRouter.GET("/feed/", controller.VerifyToken, controller.Feed)
	apiRouter.GET("/user/", controller.VerifyToken, controller.UserInfo)
	apiRouter.POST("/user/register/", controller.LoginLimit, controller.Register)
	apiRouter.POST("/user/login/", controller.LoginLimit, controller.Login)
	apiRouter.POST("/publish/action/", controller.VerifyToken, controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.VerifyToken, controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.VerifyToken, controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.VerifyToken, controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
