package controllers

import (
	"blog/Help"
	"blog/Services"
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

type Article struct {
	Id          primitive.ObjectID ` json:"id" bson:"_id"`
	Title       string             `json:"title"`
	Describe    string             `json:"describe"`
	Status      string             `json:"status"`
	Group_Photo []string           `json:"group_photo"`
	Details     string             `json:"details"`
	Disclosure  string             `json:"disclosure"`
	TopArticle  string             `json:"top_article"`
	Password    string             `json:"password"`
	Pv          int                `json:"pv"`
	//Images      string             `json:"images"`
}


type Person struct {
	ID string `uri:"id" binding:"required"`
}

func ArticleCreate(c *gin.Context) {

	article := new(Services.Article)
	article.Title = c.PostForm("title")
	article.Describe = c.PostForm("describe")
	var group_photo []string
	for _, v := range c.PostFormMap("group_photo") {
		group_photo = append(group_photo, v)
	}
	article.Group_Photo = group_photo
	article.Details = c.PostForm("details")
	article.Disclosure = c.PostForm("disclosure")
	article.Id = primitive.NewObjectID()
	article.Status = c.PostForm("status")
	article.Password = Help.Md5Encryption(c.PostForm("password"))
	Services.ServiceArticleCreate(article)
	category := c.PostFormMap("category")
	Services.ServiceArticleCreateCategory(article.Id.Hex(), category)
	Services.ServiceArticleCreateLabel(article.Id.Hex(), c.PostFormMap("label"))
	data := map[string]interface{} {"data": "创建成功"}
	Help.ReturnResponse(data, http.StatusOK, "success", c);

}

func ArticleShow(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "展示"}

	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	id, _ := primitive.ObjectIDFromHex(person.ID)
	result, _ := Services.ServiceArticleShow(id)
	if result.Disclosure == "2" {
		data["data"] = "非法操作"
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	data["data"] = result
	Help.ReturnResponse(data, http.StatusOK, "error", c)
}

func ArticleAdminShow(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "展示"}

	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	id, _ := primitive.ObjectIDFromHex(person.ID)
	result, _ := Services.ServiceArticleShow(id)
	data["data"] = result
	Help.ReturnResponse(data, http.StatusOK, "error", c)
}


func ArticlePasswordShow(c *gin.Context) {
	//var category []*Category
	id, _ := primitive.ObjectIDFromHex(c.PostForm("id"))
	password := Help.Md5Encryption(c.PostForm("password"))
	result, err := Services.ServiceArticlePasswordShow(id, password)

	data := map[string]interface{} {"data": ""}

	if err != nil {
		data["data"] = "文章不存在或密码错误"
		Help.ReturnResponse(data, 400, "err", c)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	data["data"] = result
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}

func ArticleList(c *gin.Context) {
	articles,_ :=Services.ServiceArticleList()
	data := map[string]interface{} {"data": articles}
	Help.ReturnResponse(data, 200, "success", c)
}

func ArticleEdit(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "修改成功"}

	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	id, _ := primitive.ObjectIDFromHex(person.ID)

	article := new(Services.Article)
	article.Title = c.PostForm("title")
	article.Describe = c.PostForm("describe")
	var group_photo []string
	for _, v := range c.PostFormMap("group_photo") {
		group_photo = append(group_photo, v)
	}
	article.Group_Photo = group_photo
	article.Details = c.PostForm("details")
	article.Disclosure = c.PostForm("disclosure")
	article.Status = c.PostForm("status")
	article.Password = Help.Md5Encryption(c.PostForm("password"))
	article.Id = id

	category := c.PostFormMap("category")
	Services.ServiceArticleCreateCategory(article.Id.Hex(), category)
	Services.ServiceArticleCreateLabel(article.Id.Hex(), c.PostFormMap("label"))

	err :=Services.ServiceArticleEdit(id, article)
	if err != nil {
		fmt.Println(err)
		data["data"] = "修改失败"
		Help.ReturnResponse(data, 422, "error", c)
		return
	}
	data["data"] = "修改成功"
	Help.ReturnResponse(data, 200, "error", c)

}

func ArticleDelete(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "删除成功"}
	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)

	_, err := db.Collection("article").DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, deleteerr := db.Collection("article_category").DeleteMany(ctx, bson.M{"articleid": id.Hex()})
	if deleteerr != nil {
		fmt.Println(deleteerr)
	}
	_, deletelabelerr := db.Collection("article_label").DeleteMany(ctx, bson.M{"articleid": id.Hex()})
	if deletelabelerr != nil {
		fmt.Println(deletelabelerr)
	}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}

func GetCategoryArticleList(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "获取"}
	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}
	categoryid := person.ID
	articles, _ :=Services.ServiceGetCategoryArticleList(categoryid)
	data["data"] = articles
	Help.ReturnResponse(data,  http.StatusOK, "success", c)
}

func PvUpdate(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "成功"}

	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, 400, "error", c)
		return
	}

	id, _ := primitive.ObjectIDFromHex(person.ID)
	Services.ServiceArticlePvUpadte(id)
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}
