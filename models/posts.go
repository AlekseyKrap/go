package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostItem struct {
	Title string `bson:"title"`
	Desc  string `bson:"desc"`
	Key   string `bson:"_id" json:"id,omitempty"`
}

type PostItemSlice []PostItem

func (p *PostItem) GetMongoCollectionName() string {
	return "posts"
}

var ctx = context.Background()

func (post *PostItem) Insert(db *mongo.Database) error {

	coll := db.Collection(post.GetMongoCollectionName())
	_, err := coll.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	return nil
}

func (post *PostItem) Update(db *mongo.Database) error {

	type UpdateItem struct {
		Title string `bson:"title"`
		Desc  string `bson:"desc"`
	}

	u := UpdateItem{Title: post.Title, Desc: post.Desc}

	id, err := primitive.ObjectIDFromHex(post.Key)
	if err != nil {
		return err
	}

	coll := db.Collection(post.GetMongoCollectionName())
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": id}, u)
	return err

}

func GetAllTaskItems(db *mongo.Database) (PostItemSlice, error) {

	p := PostItem{}
	coll := db.Collection(p.GetMongoCollectionName())

	cur, err := coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	posts := PostItemSlice{}
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil

}
func GetPost(db *mongo.Database, id string) (PostItem, error) {

	p := PostItem{}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return p, err
	}
	coll := db.Collection(p.GetMongoCollectionName())
	res := coll.FindOne(ctx, bson.M{"_id": ID})
	if err := res.Decode(&p); err != nil {
		return p, err
	}
	return p, nil

}
