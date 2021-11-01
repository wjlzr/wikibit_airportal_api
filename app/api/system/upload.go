package system

import (
	"encoding/base64"
	"os"
	"path"
	"wiki_bit/library/i18nresponse"
	"wiki_bit/library/ids"
	"wiki_bit/library/img"

	"github.com/gin-gonic/gin"
)

// UploadByBase64 根据base64上传图片
func UploadByBase64(c *gin.Context) {

	type Params struct {
		Bs64 string `form:"bs64" json:"bs64" binding:"required"`
	}

	var params Params
	if err := c.ShouldBindJSON(&params); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	var fileName, path1 string
	idWorker, _ := ids.NewIdWorker(10)
	code, _ := idWorker.ShortId()
	fileName = "wikibit-" + code + ".jpg"
	path1 = "public/images/"
	data, err := base64.StdEncoding.DecodeString(params.Bs64)

	// 写人本地
	f, err := os.OpenFile(path1+fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	defer f.Close()
	_, _ = f.Write(data)

	//打开文件
	fileTmp, err := os.Open(path1 + fileName)
	if err != nil {
		return
	}
	defer fileTmp.Close()

	url := img.GetImage(img.Images, fileName)
	i18nresponse.Success(c, "ok", struct {
		Url string `json:"url"`
	}{Url: url})
}

//下载
func Download(c *gin.Context) {

	filePath := c.Query("url")

	// 截取图片名称本地找静态资源里匹配
	filePath = img.ImagesRoute + string([]byte(filePath)[len(filePath)-22:])

	//打开文件
	fileTmp, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer fileTmp.Close()

	//获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filePath)

	return
}
