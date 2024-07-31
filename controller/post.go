package controller

import (
	"demo_webapp/logic"
	"demo_webapp/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	//从token中获取用户信息
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//调用业务逻辑
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, nil)

}

// PostDetailHandler 帖子详情
func PostDetailHandler(c *gin.Context) {
	//参数校验
	strID := c.Param("id")
	postID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		zap.L().Error("get post detail failed", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	//调用业务
	data, err := logic.GetPostById(postID)
	if err != nil {
		zap.L().Error("logic PostDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回状态码
	ResponseSuccess(c, data)
}

// GetPostListHandler 帖子列表
func GetPostListHandler(c *gin.Context) {
	p := &models.ParamPostList{
		Page: 1,
		Size: 5,
	}

	if err := c.ShouldBind(p); err != nil {
		ResponseError(c, CodeInvalidParma)
		return
	}

	data, err := logic.GetPostList(p.Page, p.Size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)

}

// GetPostListUpgradeHandler 按页面时间或者分数获取帖子列表
func GetPostListUpgradeHandler(c *gin.Context) {
	//1、参数校验
	p := &models.ParamPostList{
		Page:  1,
		Size:  5,
		Order: models.OrderTime,
	}
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("GetPostListUpgradeHandler with param invalid", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}
	//2、业务处理
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("logic GetPostListByOrder failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3、返回响应
	ResponseSuccess(c, data)
}
