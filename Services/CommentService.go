package Services

import (
	"blog/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Comment struct {
	Id          primitive.ObjectID ` json:"id" bson:"_id"`
	ArticleID string `json:"article_id"`
	Content  string `json:"content"`
	CreateAt int64 `json:"create_at"`
	Ip string `json:"ip"`
}
func ServiceCommentCreate(comment *Comment)  {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := db.Collection("comment").InsertOne(ctx, comment)
	if err != nil {
		log.Println(err)
	}
}
func ServiceCommentShowList(articleId string) (comments []Comment) {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var comment []Comment
	fmt.Println(articleId)
	cursor, err := db.Collection("comment").Find(ctx, bson.M{"articleid": articleId})
	if err !=nil {
		log.Println(nil)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var comments Comment
		cursor.Decode(&comments)

		comment = append(comment, comments)
	}
	return comment
}