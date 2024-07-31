package controller

import (
	"demo_webapp/logic"
	"demo_webapp/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamVote)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParma)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParma, errs.Translate(trans))
		return
	}
	//从context中获取用户user_id

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//调用业务层
	if err = logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, nil)
}
