package laboratory

import (
	"blog/Help"
	"blog/Services/laboratory"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LaboratoryConverttouppercase(c *gin.Context)  {
	 theInput := c.PostForm("input")
	 theOutput := laboratory.ServiceConverttouppercase(theInput)
	 data := map[string]interface{} {"data":theOutput}
	 fmt.Println(data)
	 Help.ReturnResponse(data, http.StatusOK, "success", c)
}
func LaboratoryConverttolowercase(c *gin.Context)  {
	theInput := c.PostForm("input")
	theOutput := laboratory.ServiceConverttolowercase(theInput)
	data := map[string]interface{} {"data":theOutput}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}