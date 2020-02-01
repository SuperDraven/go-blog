package controllers

import (
	"blog/Help"
	"blog/Services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func TimeLineCreate(c *gin.Context)  {
	data := map[string]interface{} {"data": "创建成功"}
	timeline := new(Services.TimeLine)
	timeline.Id = primitive.NewObjectID()
	timeline.Title = c.PostForm("title")
	timeline.Status = c.PostForm("status")
	timeline.CreateAT = time.Now().Unix()
	error := Services.ServiceTimeLineCreate(timeline)
	if error != nil {
		data["data"] = "创建失败"
		Help.ReturnResponse(data, http.StatusUnprocessableEntity, "error", c)
		return
	}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}

func TimeLineList(c *gin.Context)  {
	timeline := Services.ServiceTimeLineList()
	data := map[string]interface{} {"data": timeline}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}