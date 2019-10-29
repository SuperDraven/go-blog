package routers

import (
	"blog/Help"
	"blog/controllers"
	"github.com/gin-gonic/gin"
)

//func GetMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//c.Next()
//		token := c.GetHeader("Authorization")
//		_, error := Help.ParseToken(token)
//		if error != nil {
//			c.Abort()
//			c.JSON(401, gin.H{
//				"message": "error",
//				"data":    "验证失败",
//				"code":    401,
//			})
//			return
//		} else {
//			c.Next()
//		}
//
//	}
//}
func SetRouter(r *gin.Engine) {

	api := r.Group("/api")
	api.GET("/articlelist", controllers.ArticleList)
	api.GET("/labellist", controllers.LabelList)
	api.GET("/categorylist", controllers.CategoryList)
	api.GET("/article/show/:id", controllers.ArticleShow)
	api.POST("/article/passwordshow", controllers.ArticlePasswordShow)
	api.POST("/upload_img", Help.UploadImg)
	api.GET("/show_img", Help.ShowImage)
	api.PUT("/article/edit/pv/:id", controllers.PvUpdate)
	admin := api.Group("/admin")
	admin.POST("/register", controllers.Register)
	admin.POST("/login", controllers.Login)
	//admin.Use(GetMiddleware())
	//admin.GET("/site", controllers.Test)

	admin.POST("/article/create", controllers.ArticleCreate)
	admin.GET("/article/list", controllers.ArticleList)
	admin.PUT("/article/edit/:id", controllers.ArticleEdit)
	admin.GET("/article/show/:id", controllers.ArticleShow)
	admin.DELETE("/article/delete/:id", controllers.ArticleDelete)

	admin.POST("/category/create", controllers.CategoryCreate)
	admin.GET("/category/list", controllers.CategoryList)
	admin.GET("/category/treelist", controllers.CategoryTreeList)
	admin.GET("/category/show/:id", controllers.CategoryShow)

	admin.PUT("/category/edit/:id", controllers.CategoryEdit)
	admin.DELETE("/category/delete/:id", controllers.CategoryDelete)

	admin.POST("/label/create", controllers.LabelCreate)
	admin.GET("/label/list", controllers.LabelList)
	admin.GET("/label/show/:id", controllers.LabelShow)
	admin.PUT("/label/edit/:id", controllers.LabelEdit)
	admin.DELETE("/label/delete/:id", controllers.LabelDelete)

}
