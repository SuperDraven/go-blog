package Services

import (
	"blog/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
type Label struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `json:"title"`
	Color string             `json:"color"`
}

func ServiceLabelList() (labels []*Label, err error) {
	db := models.ConnectDb()
	var label []*Label
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("label").Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var labels Label
		cursor.Decode(&labels)

		label = append(label, &labels)
	}
	return label, err
}

func ServiceLabelShow(id primitive.ObjectID) (labels []*Label, errs error) {
	db := models.ConnectDb()
	var label []*Label
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("label").FindOne(ctx, bson.M{"_id": id}).Decode(&label)

	return label, err
}

func ServiceLabelCreate(label *Label) (errs error, insertErr error) {
	db := models.ConnectDb()
	filter := bson.M{"title": label.Title}

	var result struct {
		Value float64
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	error := db.Collection("label").FindOne(ctx, filter).Decode(&result)
	_, err := db.Collection("label").InsertOne(ctx, label)

	return error, err
}

func ServiceLabelEdit(id primitive.ObjectID, label *Label) (errs error) {
	db := models.ConnectDb()
	update := bson.M{"$set": label}

	_, err := db.Collection("label").UpdateOne(context.Background(), bson.M{"_id": id}, update)

	return err
}

func ServiceLabelDelete(id primitive.ObjectID) (errors error) {
	db := models.ConnectDb()
	_, err := db.Collection("label").DeleteOne(context.Background(), bson.M{"_id": id})

	return err
}