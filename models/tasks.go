package models

import (
	"database/sql"
	"fmt"
)

type PostItem struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Title      string
	Created_at string
	Updated_at string
}

type PostItemSlice []PostItem

func (post *PostItem) Insert(db *sql.DB) error {

	res, err := db.Exec(
		"INSERT INTO postsmeta SET title = ?;",
		post.Title,
	)
	id, err := res.LastInsertId()
	_, err = db.Exec("INSERT INTO postsText ( Text, idMeta) VALUES (?, ?)", post.Text, id)
	return err
}

func (post *PostItem) Update(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE postsmeta INNER JOIN poststext on postsmeta.id = poststext.idMeta SET postsmeta.title = ?, poststext.text = ? WHERE postsmeta.id = ?",
		post.Title, post.Text, post.ID,
	)
	return err
}

func GetAllTaskItems(db *sql.DB) (PostItemSlice, error) {
	rows, err := db.Query("SELECT t1.id,title,created_at,updated_at,t2.text FROM postsmeta AS t1 \nINNER JOIN poststext AS t2 on t1.id = t2.idMeta;")
	if err != nil {
		return nil, err
	}
	posts := make(PostItemSlice, 0, 8)
	for rows.Next() {
		post := PostItem{}
		if err = rows.Scan(&post.ID, &post.Title, &post.Created_at, &post.Updated_at, &post.Text); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, err
}
func GetPost(db *sql.DB, ID string) (PostItem, error) {
	post := PostItem{}
	rows, err := db.Query("SELECT  t1.id,  title,  created_at,  updated_at,  t2.text FROM postsmeta AS t1 INNER JOIN poststext AS t2 on t1.id = t2.idMeta WHERE t1.id = ?;",
		ID)
	if err != nil {
		return post, err
	}
	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.Title, &post.Created_at, &post.Updated_at, &post.Text); err != nil {
			return post, err
		}
	}
	fmt.Println(post)
	return post, err
}
