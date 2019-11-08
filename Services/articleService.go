package Services

import (
	"blog/Help"
	"blog/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
type ArticleLabelList struct {
	ArticleId string `json:"article_id"`
	LabelId   string `json:"label_id"`
}
type ArticleCategoryList struct {
	ArticleId  string `json:"article_id"`
	CategoryId string `json:"category_id"`
}
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
type Articlecategorylist []string
type Articlelabelllist []string

//创建文章
func ServiceArticleCreate(article *Article)  {
	db := models.ConnectDb()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := db.Collection("article").InsertOne(ctx, article)
	if err != nil {
		fmt.Println(err)
	}
}
//创建文章关联分类
func ServiceArticleCreateCategory(articleId string, data map[string]string) {
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
//创建文章关联标签
func ServiceArticleCreateLabel(articleId string, data map[string]string) {
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
//文章展示
func ServiceArticleShow(id primitive.ObjectID) (articles *articleShow, err error) {
	db := models.ConnectDb()
	var article *articleShow

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.Collection("article").FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	categoryIds, categoryIdsArray := ServiceArticleCategoryShow(id.Hex())
	article.Articlecategorylist = categoryIdsArray
	article.Articlelabelllist = ServiceArticleLabelShow(id.Hex())
	article.Tree = Help.GetTree(categoryIds)
	if err != nil {
		log.Fatal(err)
	}
	return article, err
}
//需要密码的文章展示
func ServiceArticlePasswordShow(id primitive.ObjectID, password string) (articles *articleShow, err error) {
	db := models.ConnectDb()
	var article *articleShow
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.Collection("article").FindOne(ctx, bson.M{"_id": id, "password": password}).Decode(&article)
	if err != nil {
		log.Println(err)
	} else {
		categoryIds, categoryIdsArray := ServiceArticleCategoryShow(id.Hex())
		article.Articlecategorylist = categoryIdsArray
		article.Articlelabelllist = ServiceArticleLabelShow(id.Hex())
		article.Tree = Help.GetTree(categoryIds)
	}
	return article, err
}
//查询文章所属分类
func ServiceArticleCategoryShow(articleId string) (map[string]string, []string) {
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
//查询文章包含标签
func ServiceArticleLabelShow(articleId string) []string {
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
//展示文章列表
func ServiceArticleList() (articles []Article, err error) {
	db := models.ConnectDb()
	var article []Article
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("article").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articles Article
		cursor.Decode(&articles)
		if articles.Status == "1" {
			articles.Status = "展示"
		} else {
			articles.Status = "隐藏"
		}
		article = append(article, articles)
	}
	return article, err
}
//修改文章
func ServiceArticleEdit(id primitive.ObjectID, article *Article) (err error)  {
	db := models.ConnectDb()
	update := bson.M{"$set": article}

	_, err = db.Collection("article").UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}
//获取指定分类的文章
func ServiceGetCategoryArticleList(categoryId string) (articles []Article ,err error)  {
	db := models.ConnectDb()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("article_category").Find(ctx, bson.M{"categoryid": categoryId})
	if err != nil {
		log.Fatal(err)
	}
	articleIds := make(map[string]string)
	var articleIdsArray []primitive.ObjectID
	for cursor.Next(ctx) {
		var article ArticleCategoryList
		cursor.Decode(&article)
		fmt.Println(article)
		articleIds[article.ArticleId] = article.ArticleId
		objectid, _:= primitive.ObjectIDFromHex(article.ArticleId)
		articleIdsArray = append(articleIdsArray, objectid)
	}
	articleCursor, errs := db.Collection("article").Find(ctx, bson.M{"_id": bson.M{"$in": articleIdsArray}})
	if errs != nil {
		fmt.Println(errs)
		log.Fatal(err)
	}
	var article []Article
	for articleCursor.Next(ctx) {
		var articles Article
		articleCursor.Decode(&articles)
		article = append(article, articles)
	}
	defer cursor.Close(ctx)
	return article, err
}
//增加文章点击
func ServiceArticlePvUpadte(id primitive.ObjectID) {
	db := models.ConnectDb()
	var article articleShow
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := db.Collection("article").FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		fmt.Println(err)
	}
	update := bson.M{"$set": bson.D{
		{"pv", article.Pv + 1},
	}}
	_, update_err := db.Collection("article").UpdateOne(ctx, bson.M{"_id": id}, update)
	if update_err != nil {
		fmt.Println(update_err)
	}
}