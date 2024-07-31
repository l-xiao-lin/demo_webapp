package logic

import (
	"demo_webapp/dao/mysql"
	"demo_webapp/models"
	"demo_webapp/pkg/snowflake"
)

// GetCommunityList 获取社区列表
func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

// Community 创建社区
func Community(p *models.ParamCommunity) error {
	//雪花算法生成community_id
	communityID := snowflake.GenID()

	//构建一个community
	community := &models.Community{
		ID:           communityID,
		Name:         p.Name,
		Introduction: p.Introduction,
	}

	//保存进数据库
	return mysql.Community(community)
}

func CommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.CommunityDetail(id)

}
