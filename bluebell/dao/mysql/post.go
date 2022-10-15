package mysql

import (
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) error {
	sqlStr := `insert into post(
		post_id,title,content,author_id,community_id)
			values(?,?,?,?,?)`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

func GetPostDetailsByID(id int64) (*models.Post, error) {
	post := new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post
	where post_id=?`
	err := db.Get(post, sqlStr, id)
	fmt.Println("post:", post)
	return post, err
}

func GetPostList(page, size int64) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 2) //这里不要写成posts := make([]*models.Post,2)
	sqlStr := `select 
		post_id,title,content,author_id,community_id,create_time
		from post
		order by create_time
		desc 
		limit ?,?`
	err := db.Select(&posts, sqlStr, (page-1)*size, size)
	return posts, err
}

//GetPostListByIDs 根据给定的id列表查询帖子的数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select
	post_id,title,content,author_id,community_id,create_time
	from post
	where post_id in (?)
	order by find_in_set(post_id,?)
`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
