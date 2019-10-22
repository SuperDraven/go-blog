package controllers

import (
	"blog/Help"
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

type Category struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `json:"title"`
	Parent_Id string             `json:"parent_id"`
}

func CategoryList(c *gin.Context) {
	db := models.ConnectDb()
	var people []*Category
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("category").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Category
		cursor.Decode(&person)

		people = append(people, &person)
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    people,
		"code":    200,
	})
}

func CategoryTreeList(c *gin.Context) {
	a := make(map[string]string)
	tree := Help.GetTree(a)
	c.JSON(200, gin.H{
		"message": "success",
		"data":    tree,
		"code":    200,
	})
}
func CategoryShow(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	var result struct {
		Id        primitive.ObjectID `bson:"_id" json:"id"`
		Title     string             `json:"title"`
		Parent_Id string             `json:"parent_id"`
	}
	//var category []*Category
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("category").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
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
func CategoryCreate(c *gin.Context) {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	category := new(Category)
	category.Title = c.PostForm("title")
	category.Parent_Id = c.PostForm("parent_id")
	category.Id = primitive.NewObjectID()
	_, err := db.Collection("category").InsertOne(ctx, category)
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

func CategoryEdit(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}

	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)
	category := new(Category)
	category.Title = c.PostForm("title")
	category.Parent_Id = c.PostForm("parent_id")
	category.Id = id
	update := bson.M{"$set": category}
	fmt.Println(id)
	_, err := db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "success",
			"data":    "修改失败",
			"code":    http.StatusUnprocessableEntity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "修改成功",
		"code":    http.StatusOK,
	})
}

func CategoryDelete(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	db := models.ConnectDb()
	SearchId := person.ID
	var result struct {
		Value float64
	}
	filter := bson.M{"parent_id": SearchId}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	error := db.Collection("category").FindOne(ctx, filter).Decode(&result)
	if error == nil {
		fmt.Println("查到了")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "success",
			"data":    "请先删除他的子分类",
			"code":    http.StatusUnprocessableEntity,
		})
	} else {
		Id, _ := primitive.ObjectIDFromHex(person.ID)
		_, err := db.Collection("category").DeleteOne(context.Background(), bson.M{"_id": Id})
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    "删除成功",
			"code":    200,
		})
	}
}
