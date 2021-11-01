package wikibit

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/fmt"
	"wiki_bit/library/openapi"
	"wiki_bit/library/random"
	"wiki_bit/services/usercenter"
)

//统一请求分发
func request(method, url string, body io.Reader) (request *http.Request, err error) {
	usercenter.PassiveInit()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Logger().Error("UserCenter request http err：", zap.Error(err))
		return request, err
	}
	req.Header.Add("Authorization", usercenter.Authorization)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

//返回参数统一处理
func responseHandle(request *http.Request) []byte {
	client := &http.Client{}

	resp, _ := client.Do(request)
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Color.Println("调用bit接口")
	fmt.Color.Printf(resp.Request)
	fmt.Color.Println("bit返回值")
	fmt.Color.Println(string(content))
	return content
}

// 活动充币
func Recharge(m RechargeRequest, communication string) (*RechargeResponseData, error) {

	m.Currency = "WKB"
	m.Nonce = random.Rand(16)
	m.Sign = sign(m)
	jsonStr, _ := json.Marshal(m)
	fmt.Color.Println("请求BIT接口前打印请求的通信唯一标识")
	fmt.Color.Println(communication)
	fmt.Color.Println("请求BIT接口参数")
	fmt.Color.Printf(m)
	request, err := request(http.MethodPost, config.Conf().Wikibit.Gateway+recharge, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Logger().Error("WikiBit Recharge 请求 err：", zap.Error(err))
		return nil, err
	}

	content := responseHandle(request)
	var result RechargeResponse
	_ = json.Unmarshal(content, &result)
	if result.Code != 200 || result.Success != true {
		return nil, err
	}
	return &result.Data, nil
}

func sign(m RechargeRequest) string {
	r := openapi.Json2UrlValues(m)
	params, _ := url.QueryUnescape(r.Encode())
	query := params + "&key=" + m.Secret
	return openapi.Sign(query)
}
