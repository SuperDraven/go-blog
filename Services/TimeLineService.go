package Services

import (
	"blog/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TimeLine struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `json:"title"`
	Status string `json:"status"`
	CreateAT int64 `json:"create_at"`
}
func ServiceTimeLineCreate(timeline *TimeLine) (errors error) {
	db :=models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := db.Collection("timeline").InsertOne(ctx, timeline)
	return err
}

func ServiceTimeLineList() (timelines []*TimeLine)  {
	db :=models.ConnectDb()
	var timeline []*TimeLine
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, _ := db.Collection("timeline").Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var timelines TimeLine
		cursor.Decode(&timelines)

		timeline = append(timeline, &timelines)
	}
	return timeline
}