package mysql

import (
	"GoAdvance/StudyGinAdvance/bluebell/models"
	"fmt"
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
		limit ?,?`
	err := db.Select(&posts, sqlStr, (page-1)*size, size)
	return posts, err
}
