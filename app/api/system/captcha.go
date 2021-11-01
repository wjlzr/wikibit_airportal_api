package system

import (
	"wiki_bit/library/captcha"
	"wiki_bit/library/i18nresponse"

	"github.com/mojocn/base64Captcha"

	"github.com/gin-gonic/gin"
)

// GenerateCaptchaHandler 获取验证码
func GenerateCaptchaHandler(c *gin.Context) {

	id, b64s, err := captcha.DriverDigitFunc()
	if err != nil {
		i18nresponse.Error(c, "500")
	}

	i18nresponse.Success(c, "ok", struct {
		B64s string `json:"b64s"`
		Uuid string `json:"uuid"`
	}{B64s: b64s, Uuid: id})
}

// CheckCode check验证码
func CheckCode(c *gin.Context) {

	type Params struct {
		Code string `form:"code" json:"code" binding:"required"`
		Uuid string `form:"uuid" json:"uuid" binding:"required"`
	}
	var params Params
	if err := c.ShouldBindJSON(&params); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	if !base64Captcha.DefaultMemStore.Verify(params.Uuid, params.Code, true) {
		i18nresponse.Error(c, "_S000023")
		return
	}

	i18nresponse.Success(c, "ok", nil)
}
