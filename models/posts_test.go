package models

import (
	"../models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("blog")

	post := models.PostItem{Desc: "11111", Title: "11111", Key: "111"}

	err := post.Insert(db)

	if err != nil {
		t.Errorf("Error Insert: %v", err)
	}

}
func TestUpdate(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("blog")

	post := models.PostItem{Desc: "11111", Title: "11111"}

	err := post.Update(db)

	if err != nil {
		t.Errorf("Error Insert: %v", err)
	}

}
func TestGetAllTaskItems(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("blog")

	_, err := models.GetAllTaskItems(db)

	if err != nil {
		t.Errorf("Error TestGetAllTaskItems: %v", err)
	}

}
func TestGetPost(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("blog")

	_, err := models.GetPost(db, "1")

	if err != nil {
		t.Errorf("Error TestGetPost: %v", err)
	}

}
