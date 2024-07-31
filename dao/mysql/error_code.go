package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("无效的密码")
	ErrorNoRow           = errors.New("没有找到相关数据")
	ErrorInvalidID       = errors.New("无效的ID")
)
