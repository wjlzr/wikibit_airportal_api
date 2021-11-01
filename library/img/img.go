package img

import "wiki_bit/config"

const (
	Images      = "images"
	ImagesRoute = "public/images/"
)

//获取图片地址
func GetImage(name, img string) string {
	return config.Conf().StaticResources.Url + img
}
