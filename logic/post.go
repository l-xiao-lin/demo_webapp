package logic

import (
	"demo_webapp/dao/mysql"
	"demo_webapp/dao/redis"
	"demo_webapp/models"
	"demo_webapp/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1、生成post_id
	p.ID = snowflake.GenID()

	//2、调用dao层
	if err = mysql.CreatePost(p); err != nil {
		return
	}
	//3、往redis中写数据
	return redis.CreatePost(p.ID, p.CommunityID)

}

// GetPostById 获取帖子详情
func GetPostById(pid int64) (postDetail *models.ApiPostDetail, err error) {
	//查询出post信息
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql GetPostById failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	//通过上面的author_id查询出authorName
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	//根据社区id查询出community
	community, err := mysql.CommunityDetail(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql CommunityDetail failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}

	postDetail = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return

}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql GetPostList failed", zap.Error(err))
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {

		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		community, err := mysql.CommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql CommunityDetail failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)

	}
	return
}

// GetPostListByOrder 通过时间或者分数获取帖子列表
func GetPostListByOrder(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//1、根据传递order参数，从redis中获取postIds
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	zap.L().Debug("redis GetPostIDsInOrder", zap.Any("ids", ids))

	//2、从mysql中查询出基于postIds的帖子
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("mysql GetPostListByIDs", zap.Any("posts", posts))

	//3、从post:vote:postID ZSet中获取每个帖子的赞成票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//4、拼接用户、社区等相关信息
	data = make([]*models.ApiPostDetail, 0, len(ids))
	for index, post := range posts {

		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			continue
		}

		community, err := mysql.CommunityDetail(post.CommunityID)
		if err != nil {
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[index],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return

}

// GetCommunityPostListByOrder 基于社区 再根据时间或者分数获取帖子列表
func GetCommunityPostListByOrder(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//1、将社区与时间、分数zset进行Zinterstore取交集，生成新的zset，并获取ids
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	//2、根据ids查询mysql
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//3、根据ids查询redis中post:vote:postID 中的投票信息
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//4、拼接用户、社区等其余信息
	data = make([]*models.ApiPostDetail, 0, len(ids))
	for index, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			continue
		}
		community, err := mysql.CommunityDetail(post.CommunityID)
		if err != nil {
			continue
		}

		detailData := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[index],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, detailData)

	}
	return
}

// GetPostListNew 两个接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostListByOrder(p)
	} else {
		data, err = GetCommunityPostListByOrder(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
	}
	return
}
