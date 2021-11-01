package openapi

import (
	"github.com/gin-gonic/gin"
	"wiki_bit/app/model/openapi/googlemap"
	"wiki_bit/library/i18nresponse"
	googlemap2 "wiki_bit/services/googlemap"
)

// 给wikiglobal调用谷歌地图经纬度中转用
func FindCoordinateByAddress(c *gin.Context) {

	var req googlemap.FindCoordinateByAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	// 请求google
	result, err := googlemap2.FindCoordinateByAddress(req.Address)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	i18nresponse.Success(c, "ok", result)
}
