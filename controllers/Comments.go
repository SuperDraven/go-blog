package controllers

import (
	"blog/Help"
	"blog/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func CommentCreate(c *gin.Context)  {
	comment := new(Services.Comment)
	comment.Id = primitive.NewObjectID()
	comment.ArticleID = c.PostForm("article_id")
	comment.Content = c.PostForm("content")
	comment.CreateAt = time.Now().Unix()
	comment.Ip = c.ClientIP()
	Services.ServiceCommentCreate(comment)
	data:=map[string]interface{} {"data":"评论成功"}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}
func CommentShowList(c *gin.Context)  {
	var person Person
	data:=map[string]interface{} {"data":""}

	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, http.StatusBadRequest, "error", c)
		return
	}
	comments := Services.ServiceCommentShowList(person.ID)
	data["data"] = comments
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}