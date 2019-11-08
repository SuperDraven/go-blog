package Services

import (
	"blog/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)
type MessageBoard struct {
	Id          primitive.ObjectID ` json:"id" bson:"_id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Ip string `json:"ip"`
	Content string `json:"content"`
	CreateAt int64 `json:"create_at"`
}
func ServiceMessageBoardCreate(messageBoard *MessageBoard)  {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := db.Collection("message_board").InsertOne(ctx, messageBoard)
	if err != nil {
		log.Println(err)
	}
}
func ServiceMessageBoardShowList() (messageBoards []MessageBoard) {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var messageBoard []MessageBoard
	cursor, err := db.Collection("message_board").Find(ctx, bson.M{})
	if err !=nil {
		log.Println(nil)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var messageBoards MessageBoard
		cursor.Decode(&messageBoards)

		messageBoard = append(messageBoard, messageBoards)
	}
	return messageBoard
}