package my_models

import (
	"../models"
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PostItem struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Title      string
	Created_at string
	Updated_at string
}

type PostItemSlice []PostItem

var ctx = context.Background()

func (post *PostItem) Insert(db *sql.DB) error {

	//res, err := db.Exec(
	//	"INSERT INTO postsmeta SET title = ?;",
	//	post.Title,
	//)
	//id, err := res.LastInsertId()
	//_, err = db.Exec("INSERT INTO postsText ( Text, idMeta) VALUES (?, ?)", post.Text, id)
	//return err

	meta := &models.Postsmetum{
		Title: post.Title,
	}
	tx := &models.Poststext{
		Text: post.Text,
	}
	err := tx.SetIdMetum(ctx, db, true, meta)
	if err != nil {
		return err
	}
	err = tx.Insert(ctx, db, boil.Infer())
	if err != nil {
		return err
	}
	return err
}

func (post *PostItem) Update(db *sql.DB) error {

	posts, err := models.Postsmeta(
		qm.Where("id=?", post.ID),
	).All(ctx, db)
	if err != nil {
		return err
	}
	_, err = posts.UpdateAll(ctx, db, models.M{"Title": post.Title})
	if err != nil {
		return err
	}

	postsT, err := models.Poststexts(
		qm.Where("idMeta=?", post.ID),
	).All(ctx, db)
	if err != nil {
		return err
	}
	_, err = postsT.UpdateAll(ctx, db, models.M{"Text": post.Text})
	if err != nil {
		return err
	}
	return err

}

func GetAllTaskItems(db *sql.DB) (PostItemSlice, error) {
	rows, err := models.NewQuery(
		qm.Select("t1.id", "t1.title", "t1.created_at", "t1.updated_at", "t2.text"),
		qm.From("postsmeta AS t1"),
		qm.InnerJoin("poststext AS t2 on t1.id = t2.idMeta"),
	).Query(db)
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
	rows, err := models.NewQuery(
		qm.Select("t1.id", "t1.title", "t1.created_at", "t1.updated_at", "t2.text"),
		qm.From("postsmeta AS t1"),
		qm.InnerJoin("poststext AS t2 on t1.id = t2.idMeta"),
		qm.Where("t1.id = ?", ID),
	).Query(db)
	if err != nil {
		return post, err
	}
	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.Title, &post.Created_at, &post.Updated_at, &post.Text); err != nil {
			return post, err
		}
	}
	return post, err

}
