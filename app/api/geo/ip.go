package geo

import (
	"wiki_bit/library/geo"
	"wiki_bit/library/response"

	"github.com/gin-gonic/gin"
)

//GetWithIpToLocation 根据ip返回国家
func GetWithIpToLocation(c *gin.Context) {

	code := geo.GetCountryCode(c.ClientIP())

	response.Success(c, struct {
		Code string `json:"code"`
	}{Code: code}, "ok")
}
