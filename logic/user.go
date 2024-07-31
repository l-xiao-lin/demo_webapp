package logic

import (
	"demo_webapp/dao/mysql"
	"demo_webapp/models"
	"demo_webapp/pkg/jwt"
	"demo_webapp/pkg/snowflake"
)

// SignUp 注册
func SignUp(p *models.ParamSignUp) (err error) {
	//1、判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//2、生成uuid
	userID := snowflake.GenID()

	//3、构造User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//4、保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	//查询数据库中的是否有该用户并且密码是否一致
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return
	}

	//生成token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return

}
