package mysql

import (
	"database/sql"
	"demo_webapp/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(pid int64) (postDetail *models.Post, err error) {
	postDetail = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,status,create_time from post where post_id = ?`
	if err = db.Get(postDetail, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return
}

func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,status,create_time from post order by create_time desc limit ?,?`
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	return

}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,status,create_time from post where post_id in (?) order by find_in_set (post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return

}
