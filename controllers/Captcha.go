package controllers

import (
	"blog/Help"
	"blog/Services"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
)


func GenerateCaptchaHandler(c *gin.Context)  {
	captcha := Services.ServiceGenerateCaptchaHandler()
	data := map[string]interface{}{"data": captcha}
	Help.ReturnResponse(data, http.StatusOK, "success", c)
}
func Getcaptcha(c *gin.Context) {
	//captchaId := c.Param("captchaId")
	//fmt.Println("GetCaptchaPng : " + captchaId)
	Services.ServeHTTP(c.Writer, c.Request)
}
func CaptchaVerify(c *gin.Context) {
	captchaId := c.Param("captchaId")
	value := c.Param("value")
	data := map[string]interface{}{"data": "验证码错误"}
	if captchaId == "" || value == "" {
		data["data"] = "参数错误"
		Help.ReturnResponse(data, http.StatusUnprocessableEntity, "error", c)
		return
	}
	if captcha.VerifyString(captchaId, value) {
		data["data"] = "验证成功"
		Help.ReturnResponse(data, http.StatusOK, "success", c)
		return
	}
	Help.ReturnResponse(data, http.StatusUnprocessableEntity, "error", c)
}
