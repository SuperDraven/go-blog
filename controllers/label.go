package controllers

import (
	"blog/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type Label struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `json:"title"`
	Color string             `json:"color"`
}

func LabelList(c *gin.Context) {
	db := models.ConnectDb()
	var label []*Label
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("label").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Label
		cursor.Decode(&person)

		label = append(label, &person)
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    label,
		"code":    200,
	})

}
func LabelShow(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	var result struct {
		Id    primitive.ObjectID `bson:"_id" json:"id"`
		Title string             `json:"title"`
		Color string             `json:"color"`
	}
	//var category []*Category
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("label").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	fmt.Println(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
		"code":    http.StatusOK,
	})
}
func LabelCreate(c *gin.Context) {
	db := models.ConnectDb()
	label := new(Label)
	label.Title = c.PostForm("title")
	label.Color = c.PostForm("color")
	label.Id = primitive.NewObjectID()

	var result struct {
		Value float64
	}
	filter := bson.M{"title": label.Title}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	error := db.Collection("label").FindOne(ctx, filter).Decode(&result)
	if error == nil {
		fmt.Println("查到了")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "success",
			"data":    "标签标题重复",
			"code":    http.StatusUnprocessableEntity,
		})
	} else {
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		_, err := db.Collection("label").InsertOne(ctx, label)
		if err != nil {
			fmt.Println(err)
			c.JSON(422, gin.H{
				"message": "success",
				"data":    "创建失败",
				"code":    200,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "success",
				"data":    "创建成功",
				"code":    200,
			})
		}
	}

}

func LabelEdit(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)
	label := new(Label)
	label.Title = c.PostForm("title")
	label.Color = c.PostForm("color")
	label.Id = id
	update := bson.M{"$set": label}
	fmt.Println(id)
	_, err := db.Collection("label").UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		fmt.Println(err)
		c.JSON(422, gin.H{
			"message": "success",
			"data":    "修改失败",
			"code":    200,
		})
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    "修改成功",
		"code":    200,
	})
}

func LabelDelete(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)

	_, err := db.Collection("label").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "删除成功",
		"code":    200,
	})
}
