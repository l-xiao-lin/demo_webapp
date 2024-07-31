package mysql

import (
	"database/sql"
	"demo_webapp/models"
	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id,community_name,introduction from community`

	if err = db.Select(&communityList, sqlStr); err != nil {
		return nil, err
	}

	if len(communityList) == 0 {
		zap.L().Warn("no community found")
		return nil, ErrorNoRow
	}
	return
}

// Community 创建社区
func Community(community *models.Community) (err error) {
	sqlStr := `insert into community(community_id,community_name,introduction) values(?,?,?)`
	_, err = db.Exec(sqlStr, community.ID, community.Name, community.Introduction)
	return
}

// CommunityDetail 社区详情
func CommunityDetail(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id=?`
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}

	return
}
