package configure

import (
	"github.com/gin-gonic/gin"
	"wiki_bit/app/model/configure"
	"wiki_bit/library/i18nresponse"
)

func GetConfigure(c *gin.Context) {

	var m configure.Configure

	m, err := m.FindOne()
	if err != nil {
		i18nresponse.Error(c, "500")
	}

	i18nresponse.Success(c, "ok", m)
}
