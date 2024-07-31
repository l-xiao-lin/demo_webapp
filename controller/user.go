package controller

import (
	"demo_webapp/dao/mysql"
	"demo_webapp/logic"
	"demo_webapp/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 注册
func SignUpHandler(c *gin.Context) {
	//1、参数处理
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParma)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParma, errs.Translate(trans))

		return
	}

	//2、业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return

	}

	//3、返回响应

	ResponseSuccess(c, nil)

}

func LoginHandler(c *gin.Context) {
	//1、参数处理
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParma)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParma, errs.Translate(trans))
		return
	}

	//2、业务处理
	data, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3、返回响应

	ResponseSuccess(c, gin.H{
		"user_id":  data.UserID,
		"username": data.Username,
		"token":    data.Token,
	})
}
