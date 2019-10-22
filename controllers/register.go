package controllers

import (
	"blog/Help"
	"blog/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func Register(c *gin.Context)  {
	db :=models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	name :=c.PostForm("name")
	password :=c.PostForm("password")
	password = Help.Md5Encryption(password)
	res, _ := db.Collection("admin_users").InsertOne(ctx, bson.M{"name": name, "password": password})
	id := res.InsertedID
	user :=new(Help.UserInfo)
	user.Name = name
	user.Password = password
	token, _ := Help.CreateToken(user)


	fmt.Println(user)
	fmt.Println(id)
	fmt.Println(token)
	c.JSON(200, gin.H{
		"data":"bearer " + token,
	})
}