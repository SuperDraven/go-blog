package controllers

import (
	"blog/Help"
	"blog/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var User struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string
	//Password string
}

func Login(c *gin.Context) {
	db := models.ConnectDb()
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	name := c.PostForm("name")
	password := c.PostForm("password")
	password = Help.Md5Encryption(password)

	//res, _ := db.Collection("admin_users").InsertOne(ctx, bson.M{"name": name, "password": password})
	//id := res.InsertedID
	user := new(Help.UserInfo)
	user.Name = name
	user.Password = password
	fmt.Println(user)
	//filter := bson.M{"name": "pi"}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("admin_users").FindOne(ctx, user).Decode(&User)
	if err != nil {
		fmt.Println(err)
		c.JSON(422, gin.H{
			"message": "error",
			"data":    "用户名密码错误",
			"code":    422,
		})
	} else {
		token, _ := Help.CreateToken(user)
		c.JSON(200, gin.H{
			"message": "success",
			"data":    "bearer " + token,
			"code":    200,
		})
	}
	//fmt.Println(user)
	//fmt.Println(id)
	//fmt.Println(token)

}
