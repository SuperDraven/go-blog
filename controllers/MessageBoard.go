package controllers

import (
	"blog/Help"
	"blog/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func MessageBoardCreate(c *gin.Context)  {
	messageBoard := new(Services.MessageBoard)
	messageBoard.Ip = c.ClientIP()
	messageBoard.Id = primitive.NewObjectID()
	messageBoard.CreateAt = time.Now().Unix()
	messageBoard.Email = c.PostForm("email")
	messageBoard.Name = c.PostForm("name")
	messageBoard.Content = c.PostForm("content")
	Services.ServiceMessageBoardCreate(messageBoard)
	data := map[string]interface{} {"data":"留言成功"}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}

func MessageBoardShowList(c *gin.Context)  {
	messageBoardList := Services.ServiceMessageBoardShowList()
	data := map[string]interface{} {"data":messageBoardList}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}