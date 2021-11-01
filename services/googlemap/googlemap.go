package googlemap

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/services"

	"github.com/buger/jsonparser"
	"go.uber.org/zap"
)

const getWay = "https://maps.googleapis.com/maps/api/geocode/json"

type Coordinate struct {
	Lat float64 `json:"lat"` // 维度
	Lng float64 `json:"lng"` // 经度
}

// FindCoordinateByAddress
func FindCoordinateByAddress(address string) (coordinate Coordinate, err error) {

	// 地址为空则直接返回
	if address == "" {
		return coordinate, errors.New("1010003")
	}

	response, err := services.Request(http.MethodGet, getWay+"?address="+url.QueryEscape(address)+"&key="+config.Conf().GoogleMap.Key, nil)
	if err != nil {
		log.Logger().Error("地址转经纬度 err: ", zap.Error(err))
		return coordinate, err
	}
	content, err := services.ResponseHandle(response)
	if err != nil {
		return coordinate, err
	}
	status, err := jsonparser.GetString(content, "status")
	if err != nil {
		log.Logger().Error("jsonparser 解析错误 err: ", zap.Error(err))
		return coordinate, err
	}

	if status != "OK" {
		errorMessage, _ := jsonparser.GetString(content, "error_message")
		log.Logger().Error("googlemap 调用地理编码失败 err: ", zap.Error(errors.New(errorMessage)))
		return Coordinate{}, errors.New("googlemap 调用地理编码失败")
	}

	value, _, _, err := jsonparser.Get(content, "results", "[0]", "geometry", "location")
	if err != nil {
		log.Logger().Error("googlemap 地理编码转换失败 err: ", zap.Error(err))
		return Coordinate{}, errors.New("googlemap 地理编码转换失败")
	}
	_ = json.Unmarshal(value, &coordinate)
	return
}
