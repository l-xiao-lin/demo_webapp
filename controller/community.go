package controller

import (
	"demo_webapp/logic"
	"demo_webapp/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityListHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CommunityHandler(c *gin.Context) {
	p := new(models.ParamCommunity)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create community with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}
	if err := logic.Community(p); err != nil {
		zap.L().Error("logic Community failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func CommunityDetailHandler(c *gin.Context) {
	//参数处理
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		zap.L().Error("CommunityDetailHandler invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	//业务处理
	data, err := logic.CommunityDetail(id)
	if err != nil {
		zap.L().Error("logic CommunityDetail failed ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)

}
