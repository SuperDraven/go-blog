package Help

import "github.com/gin-gonic/gin"

func ReturnResponse(data map[string]interface{}, code int,message string ,c *gin.Context)  {
	datas := map[string]interface{} {"message": message, "code": code, "data": data["data"]}

	c.JSON(code, datas)
}

