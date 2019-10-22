package models

import (
	"blog/conf"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type User struct {
	name string
	role int
}
var db *mongo.Database

func ConnectDb() *mongo.Database {
	ctx, _:=context.WithTimeout(context.Background(), 10*time.Second)
	client, err :=mongo.Connect(ctx, options.Client().ApplyURI(conf.LoadConf().DB))
	if err != nil {
		log.Println(err)
	}
	return client.Database("blog")
}