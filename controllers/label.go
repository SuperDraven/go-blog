package controllers

import (
	"blog/Help"
	"blog/Services"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)


func LabelList(c *gin.Context) {

	label, err :=Services.ServiceLabelList()
	if err!=nil {
		log.Println(err)
	}
	data := map[string]interface{} {"data":label}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}
func LabelShow(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data":""}
	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, http.StatusBadRequest, "error", c)
		return
	}
	id, _ := primitive.ObjectIDFromHex(person.ID)
	label, err:=Services.ServiceLabelShow(id)
	if err != nil {
		log.Println(err)
	}
	data["data"] = label
	Help.ReturnResponse(data, http.StatusOK, "success", c)

}
func LabelCreate(c *gin.Context) {

	label := new(Services.Label)
	label.Title = c.PostForm("title")
	label.Color = c.PostForm("color")
	label.Id = primitive.NewObjectID()
	data := map[string]interface{} {"data": "创建成功"}
	err,insertErr := Services.ServiceLabelCreate(label)
	if err == nil {
		fmt.Println("查到了")
		data["data"] = "标题重复"
		Help.ReturnResponse(data,http.StatusUnprocessableEntity, "error", c)
		return
	}

	if insertErr != nil {
		data["data"] = "创建失败"
		Help.ReturnResponse(data,http.StatusUnprocessableEntity, "error", c)
		return
	}
	Help.ReturnResponse(data,http.StatusOK, "success", c)
}

func LabelEdit(c *gin.Context) {
	var person Person
	data := map[string]interface{} {"data": "修改成功"}
	if err := c.ShouldBindUri(&person); err != nil {
		data["data"] = err
		Help.ReturnResponse(data, http.StatusBadRequest, "error", c)
		return
	}
	id, _ := primitive.ObjectIDFromHex(person.ID)
	label := new(Services.Label)
	label.Title = c.PostForm("title")
	label.Color = c.PostForm("color")
	label.Id = id
	err := Services.ServiceLabelEdit(id, label)
	if err != nil {
		data["data"] = "修改失败"
		Help.ReturnResponse(data, http.StatusUnprocessableEntity, "error", c)
		return
	}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}

func LabelDelete(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"message": "error", "data": err, "code": 400})
		return
	}
	data := map[string]interface{} {"data":"删除成功"}
	id, _ := primitive.ObjectIDFromHex(person.ID)
	err:= Services.ServiceLabelDelete(id)
	if err != nil {
		data["data"] = "删除失败"
		Help.ReturnResponse(data,http.StatusUnprocessableEntity, "error", c)
		return
	}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}
