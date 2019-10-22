package Help

import (
	"blog/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ShowImage(c *gin.Context) {
	imageName := c.Query("imageName")
	fmt.Println(imageName)
	c.File("./img/" + imageName)
}
func UploadImg(c *gin.Context) {
	file, _ := c.FormFile("file")
	//fmt.Println(file)
	log.Println(file.Filename)
	url := conf.LoadConf().SiteUrl + ":" + conf.LoadConf().SitePort + "/api/show_img?imageName="
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./img/"+file.Filename)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    url + file.Filename,
		"code":    200,
	})

	//c.String(http.StatusOK, fmt.Sprintf(c.Request.URL.Path))
}
