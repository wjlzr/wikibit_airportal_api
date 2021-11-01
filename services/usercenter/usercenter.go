package usercenter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"wiki_bit/app/model/oauth"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/response"

	"go.uber.org/zap"
)

//init token
func PassiveInit() {
	resp, err := http.Get(config.Conf().UserCenter.SignUrl + fmt.Sprintf(getToken+"?username=%s&password=%s", userName, password))
	if err != nil {
		log.Logger().Error("UserCenter init http err：", zap.Error(err))
		response.Error(c, -1, err.Error())
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger().Error("UserCenter init ioutil.ReadAll err：", zap.Error(err))
		response.Error(c, -1, err.Error())
	}
	var t tokenResponse
	_ = json.Unmarshal(bs, &t)
	if t.Status == false || t.AccessToken == "" {
		log.Logger().Info("UserCenter init 授权失败 err：", zap.Error(err))
		response.Error(c, -1, "授权失败")
	}
	Authorization = t.TokenType + " " + t.AccessToken
}

//统一请求分发
func request(method, url string, body io.Reader) (request *http.Request, err error) {
	PassiveInit()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Logger().Error("UserCenter request http err：", zap.Error(err))
		return request, err
	}
	req.Header.Add("Authorization", Authorization)
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
	fmt.Printf("\n调用中台接口：%+v \n", resp.Request)
	fmt.Printf("\n用户中台返回值：%s \n", string(content))
	return content
}

// 发送验证码 特供bit
func SendCode(areaCode, phone, languageCode string) (s SendCodeResponse, err error) {

	request, err := request(http.MethodGet, config.Conf().UserCenter.User+fmt.Sprintf(sendCodeBit+"?mobile=%s&areaCode=%s&language=%s", phone, areaCode, languageCode), nil)
	if err != nil {
		log.Logger().Error("UserCenter SendCode 请求 err：", zap.Error(err))
		return SendCodeResponse{}, err
	}

	content := responseHandle(request)
	_ = json.Unmarshal(content, &s)
	if s.Code != 200 || s.Success != true {
		log.Logger().Info("UserCenter SendCode 发送验证码Error response：", zap.Any("response", s))
		return SendCodeResponse{}, errors.New(s.Msg)
	}
	return s, nil
}

//验证短信验证码 特供bit
func ValidateCode(req oauth.ValidateCodeRequest) (res CurrencyResponse, err error) {

	request, err := request(http.MethodGet, config.Conf().UserCenter.User+fmt.Sprintf(getCodesByRequestIdBit+"?areaCode=%s&mobile=%s&code=%s", req.AreaCode, req.PhoneNumber, req.Smscode), nil)
	if err != nil {
		log.Logger().Error("UserCenter ValidateCode 请求 err：", zap.Error(err))
		return CurrencyResponse{}, err
	}

	content := responseHandle(request)
	_ = json.Unmarshal(content, &res)
	return res, nil
}

// 发送验证码 中台
//func SendCode(areaCode, phone, languageCode, userId string, smsBusinessType int) (s SendCodeResponse, err error) {
//
//	jsonStr, _ := json.Marshal(sendCodeRequest{AreaCode: areaCode, Phone: phone, LanguageCode: languageCode, UserId: userId, SmsBusinessType: smsBusinessType, ApplicationType: applicationType})
//	request, err := request(http.MethodPost, config.Conf().UserCenter.User+sendCode, bytes.NewBuffer(jsonStr))
//	if err != nil {
//		log.Logger().Error("UserCenter SendCode 请求 err：", zap.Error(err))
//		return SendCodeResponse{}, err
//	}
//
//	content := responseHandle(request)
//	_ = json.Unmarshal(content, &s)
//	if s.Code != 200 || s.Success != true {
//		log.Logger().Info("UserCenter SendCode 发送验证码Error response：", zap.Any("response", s))
//		return SendCodeResponse{}, errors.New(s.Msg)
//	}
//	return s, nil
//}

// 验证短信验证码 用户中台
//func ValidateCode(req oauth.ValidateCodeRequest) (res CurrencyResponse, err error) {
//
//	request, err := request(http.MethodGet, config.Conf().UserCenter.User+fmt.Sprintf(validateCode+"?areaCode=%s&phoneNumber=%s&smscode=%s&userId=%s&applicationType=%d", req.AreaCode, req.PhoneNumber, req.Smscode, req.UserId, applicationType), nil)
//	if err != nil {
//		log.Logger().Error("UserCenter ValidateCode 请求 err：", zap.Error(err))
//		return CurrencyResponse{}, err
//	}
//
//	content := responseHandle(request)
//	_ = json.Unmarshal(content, &res)
//	return res, nil
//}

// 发送邮箱验证码
func SendEmailCode(m oauth.SendEmailCodeRequest) (result CurrencyResponse, err error) {
	jsonStr, _ := json.Marshal(m)
	request, err := request(http.MethodPost, config.Conf().UserCenter.User+sendEmail, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Logger().Error("UserCenter SendEmailCode 请求 err：", zap.Error(err))
		return CurrencyResponse{}, err
	}

	content := responseHandle(request)
	var resp CurrencyResponse
	_ = json.Unmarshal(content, &resp)
	if resp.Code != 200 || resp.Success != true {
		log.Logger().Info("UserCenter SendEmailCode 发送邮箱验证码Err response：", zap.Any("response", resp))
		return CurrencyResponse{}, errors.New(resp.Msg)
	}
	return resp, nil
}

// 验证邮箱验证码
func ValidateEmailCode(email, userId, code string) (result SendCodeResponse, err error) {

	request, err := request(http.MethodGet, config.Conf().UserCenter.User+fmt.Sprintf(validateEmailCode+"?email=%s&userId=%s&code=%s&applicationType=%d", email, userId, code, applicationType), nil)
	if err != nil {
		log.Logger().Error("UserCenter ValidateEmailCode 请求 err：", zap.Error(err))
		return SendCodeResponse{}, err
	}

	content := responseHandle(request)

	var resp SendCodeResponse
	_ = json.Unmarshal(content, &resp)
	if resp.Code != 200 || resp.Success != true {
		log.Logger().Info("UserCenter ValidateEmailCode 验证邮箱验证码 Err response：", zap.Any("response", resp))
		return SendCodeResponse{}, errors.New(resp.Msg)
	}

	return resp, nil
}

func UserInfo(req oauth.GetUserInfoRequest) (*CurrencyWithUserResponseData, error) {

	request, err := request(http.MethodGet, config.Conf().UserCenter.User+fmt.Sprintf(getUser+"?userId=%s&countryCode=%s&applicationType=%d", req.UserId, req.CountryCode, applicationType), nil)
	if err != nil {
		log.Logger().Error("UserCenter UserInfo 请求 err：", zap.Error(err))
		return nil, err
	}

	content := responseHandle(request)
	var res CurrencyWithUserResponse
	_ = json.Unmarshal(content, &res)
	if res.Code != 200 || res.Success != true {
		return nil, errors.New(res.Msg)
	}
	return &res.Data, nil
}
