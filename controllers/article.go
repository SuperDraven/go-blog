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
type ArticleCategoryList struct {
	ArticleId  string `json:"article_id"`
	CategoryId string `json:"category_id"`
}
type ArticleLabelList struct {
	ArticleId string `json:"article_id"`
	LabelId   string `json:"label_id"`
}
type Person struct {
	ID string `uri:"id" binding:"required"`
}
type Articlecategorylist []string
type Articlelabelllist []string
type articleShow struct {
	Id          primitive.ObjectID ` json:"id" bson:"_id"`
	Title       string             `json:"title"`
	Describe    string             `json:"describe"`
	Status      string             `json:"status"`
	Group_Photo []string           `json:"group_photo"`
	Details     string             `json:"details"`
	Disclosure  string             `json:"disclosure"`
	TopArticle  string             `json:"top_article"`
	Pv          int                `json:"pv"`

	Articlelabelllist   `json:"label"`
	Articlecategorylist `json:"category"`
	Tree                []*Help.CategoryTree `json:"tree"`
}

func ArticleCreate(c *gin.Context) {
	db := models.ConnectDb()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	article := new(Article)
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

	_, err := db.Collection("article").InsertOne(ctx, article)
	category := c.PostFormMap("category")
	createCategory(article.Id.Hex(), category)
	createLabel(article.Id.Hex(), c.PostFormMap("label"))
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
func ArticleShow(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	var result articleShow
	//var category []*Category
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("article").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	categoryIds, categoryIdsArray := articlecategoryShow(person.ID)
	//category_list := articlecategoryShow(person.ID)
	if result.Disclosure == "2" {
		c.JSON(400, gin.H{
			"message": "error",
			"data":    "非法操作",
			"code":    400,
		})
	} else {
		result.Articlecategorylist = categoryIdsArray
		result.Articlelabelllist = articlelabelShow(person.ID)
		result.Tree = Help.GetTree(categoryIds)
		//fmt.Println(category_list)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    result,
			"code":    http.StatusOK,
		})
	}

}
func ArticlePasswordShow(c *gin.Context) {
	var result articleShow
	//var category []*Category
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(c.PostForm("id"))
	password := Help.Md5Encryption(c.PostForm("password"))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("article").FindOne(ctx, bson.M{"_id": id, "password": password}).Decode(&result)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "err",
			"data":    "文章不存在或密码错误",
			"code":    400,
		})
	} else {
		categoryIds, categoryIdsArray := articlecategoryShow(c.PostForm("id"))
		//category_list := articlecategoryShow(person.ID)
		result.Articlecategorylist = categoryIdsArray
		result.Articlelabelllist = articlelabelShow(c.PostForm("id"))
		result.Tree = Help.GetTree(categoryIds)
		//fmt.Println(category_list)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    result,
			"code":    http.StatusOK,
		})
	}

}

func ArticleList(c *gin.Context) {
	db := models.ConnectDb()
	var people []Article
	//filter := bson.M{"name": "pi"}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("article").Find(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Article
		cursor.Decode(&person)
		if person.Status == "1" {
			person.Status = "展示"
		} else {
			person.Status = "隐藏"
		}
		people = append(people, person)
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    people,
		"code":    200,
	})
}
func ArticleEdit(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)

	article := new(Article)
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
	createCategory(article.Id.Hex(), category)
	createLabel(article.Id.Hex(), c.PostFormMap("label"))
	update := bson.M{"$set": article}

	_, err := db.Collection("article").UpdateOne(context.Background(), bson.M{"_id": id}, update)
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
func ArticleDelete(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
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
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "删除成功",
		"code":    200,
	})
}
func createCategory(articleId string, data map[string]string) {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, deleteerr := db.Collection("article_category").DeleteMany(ctx, bson.M{"articleid": articleId})
	if deleteerr != nil {
		fmt.Println(deleteerr)
	}
	category := new(ArticleCategoryList)
	category.ArticleId = articleId
	for _, v := range data {
		category.CategoryId = v
		_, err := db.Collection("article_category").InsertOne(ctx, category)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("article_id" + articleId + "category_id" + v)
	}
}
func createLabel(articleId string, data map[string]string) {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, deleteerr := db.Collection("article_label").DeleteMany(ctx, bson.M{"articleid": articleId})
	if deleteerr != nil {
		fmt.Println(deleteerr)
	}
	label := new(ArticleLabelList)
	label.ArticleId = articleId
	for _, v := range data {
		label.LabelId = v
		_, err := db.Collection("article_label").InsertOne(ctx, label)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//分类
func articlecategoryShow(articleId string) (map[string]string, []string) {
	db := models.ConnectDb()
	//a := [...]string{"5d9af12c04c971ca5a5e8e1f", "5d9af13104c971ca5a5e8e20"}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("article_category").Find(ctx, bson.M{"articleid": articleId})

	if err != nil {
		log.Fatal(err)
	}
	categoryIds := make(map[string]string)
	var categoryIdsArray []string
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person ArticleCategoryList
		cursor.Decode(&person)
		categoryIds[person.CategoryId] = person.CategoryId
		categoryIdsArray = append(categoryIdsArray, person.CategoryId)
	}
	fmt.Println(categoryIds)
	//a := [...]string{"5d9af12c04c971ca5a5e8e1f", "5d9af13104c971ca5a5e8e20"}
	return categoryIds, categoryIdsArray
}

//标签
func articlelabelShow(articleId string) []string {
	db := models.ConnectDb()
	//a := [...]string{"5d9af12c04c971ca5a5e8e1f", "5d9af13104c971ca5a5e8e20"}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("article_label").Find(ctx, bson.M{"articleid": articleId})

	if err != nil {
		log.Fatal(err)
	}
	var labelIds []string
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person ArticleLabelList
		cursor.Decode(&person)
		//categoryIds["aaa"] = person.CategoryId
		labelIds = append(labelIds, person.LabelId)
	}
	fmt.Println(labelIds)
	//a := [...]string{"5d9af12c04c971ca5a5e8e1f", "5d9af13104c971ca5a5e8e20"}
	return labelIds
}
func PvUpdate(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	var result articleShow
	db := models.ConnectDb()
	id, _ := primitive.ObjectIDFromHex(person.ID)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("article").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	update := bson.M{"$set": bson.D{
		{"pv", result.Pv + 1},
	}}
	_, update_err := db.Collection("article").UpdateOne(ctx, bson.M{"_id": id}, update)
	if update_err != nil {
		fmt.Println(update_err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "成功",
		"code":    http.StatusOK,
	})
}
