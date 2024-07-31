package router

import (
	"demo_webapp/controller"
	"demo_webapp/logger"
	"demo_webapp/middlewares"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	//使用第三方ginzap来接收gin框架的系统日志
	r.Use(ginzap.Ginzap(logger.LG, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.LG, true))

	//注册

	v1 := r.Group("/api/v1")

	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	//社区

	v1.GET("/community", controller.CommunityListHandler)
	v1.POST("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	//帖子
	v1.GET("/post/:id", controller.PostDetailHandler)
	v1.GET("/post", controller.GetPostListHandler)

	//按页面时间或分数获取帖子列表
	v1.GET("/postList", controller.GetPostListUpgradeHandler)

	//中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.POST("/post", controller.CreatePostHandler)

		//投票
		v1.POST("/vote", controller.PostVoteHandler)

	}
	return r
}
