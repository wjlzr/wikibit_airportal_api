package oauth

import (
	"encoding/json"
	"strings"
	"wiki_bit/app/model/oauth"
	"wiki_bit/app/model/participant"
	"wiki_bit/app/service/configure"
	oauth2 "wiki_bit/app/service/oauth"
	"wiki_bit/app/service/participate"
	participants "wiki_bit/app/service/participate"
	"wiki_bit/library/array"
	"wiki_bit/library/constant"
	"wiki_bit/library/fmt"
	"wiki_bit/library/i18nresponse"
	"wiki_bit/library/validate"
	"wiki_bit/services/usercenter"

	"github.com/gin-gonic/gin"
)

//发送短信验证码
func SmsSend(c *gin.Context) {

	code := c.Request.FormValue("code")
	phone := c.Request.FormValue("phone")
	languageCode := c.Request.FormValue("languageCode")

	// 发送验证码
	smsResult, err := usercenter.SendCode(code, phone, languageCode)
	if err != nil || smsResult.Data.Succeed != true {
		if smsResult.Data.Message == "" {
			i18nresponse.Error(c, "1010005")
		} else {
			i18nresponse.Error(c, smsResult.Data.Message)
		}
		return
	}

	i18nresponse.Success(c, "ok", struct {
		Success bool `json:"success"`
	}{Success: true})
}

// 验证短信验证码
func ValidateCode(c *gin.Context) {

	var p participant.Participant
	var req oauth.ValidateCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	// 先做人机验证
	//val, err := googlerecaptcha.CheckRecaptcha(req.RecaptchaResponse, req.PhoneNumber)
	//if err != nil {
	//	i18nresponse.Error(c, "500")
	//	return
	//}
	//if !val {
	//	i18nresponse.Error(c, "1010019")
	//	return
	//}

	result, err := usercenter.ValidateCode(req)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	if result.Success == false || result.Data.Succeed == false {
		i18nresponse.Error(c, result.Data.Message)
		return
	}

	// 短信验证成功直接走创建
	stu, err := json.Marshal(req)
	if err = json.Unmarshal(stu, &p); err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	p.EthAddress = strings.Trim(p.EthAddress, " ")

	if !validate.ValidateAddress(p.EthAddress) {
		i18nresponse.Error(c, "1010001")
		return
	}

	created, err := participants.Create(c, p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}
	participate.Logging(created)

	i18nresponse.Success(c, "ok", created)
}

// 发送邮箱验证码
//func SendEmailCode(c *gin.Context) {
//
//	var req oauth.SendEmailCodeRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.Error(c, -1, "缺少必要参数")
//		return
//	}
//
//	result, err := usercenter.SendEmailCode(req)
//	if err != nil || result.Data.Succeed != true || result.Data.Message != "success" {
//		response.Error(c, -1, "发送验证码失败")
//		return
//	}
//
//	response.Success(c, result.Data, "ok")
//}

// 发送邮箱验证码
func SendEmailCode(c *gin.Context) {

	result, err := configure.GetInfo()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if result.EnableMailbox != 1 {
		i18nresponse.Error(c, "1010020")
		return
	}

	var req oauth.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}
	req.Email = strings.Trim(req.Email, " ")
	arr := strings.Split(req.Email, "@")
	if val := array.StrInArray(constant.Emails, arr[1]); !val {
		i18nresponse.Error(c, "_S000032")
		return
	}

	if _, err := oauth2.SendEmailCode(req); err != nil {
		i18nresponse.Error(c, "1010002")
		return
	}

	i18nresponse.Success(c, "ok", nil)
}

// 验证邮箱验证码
func ValidateEmailCode(c *gin.Context) {

	result, err := configure.GetInfo()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if result.EnableMailbox != 1 {
		i18nresponse.Error(c, "1010020")
		return
	}

	var p participant.Participant
	var req oauth.ValidateEmailCode
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	// 先做人机验证
	//val, err := googlerecaptcha.CheckRecaptcha(req.RecaptchaResponse, req.Communication)
	//if err != nil {
	//	i18nresponse.Error(c, "500")
	//	return
	//}
	//if !val {
	//	i18nresponse.Error(c, "1010019")
	//	return
	//}
	req.Communication = strings.Trim(req.Communication, " ")
	// 过滤国内邮箱
	arr := strings.Split(req.Communication, "@")
	if val := array.StrInArray(constant.Emails, arr[1]); !val {
		i18nresponse.Error(c, "_S000032")
		return
	}

	if val, err := oauth2.ValidateEmailCode(req.Communication, req.Code); err != nil || !val {
		i18nresponse.Error(c, "1010004")
		return
	}

	// 邮箱验证成功直接走创建
	stu, err := json.Marshal(req)
	if err = json.Unmarshal(stu, &p); err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	p.EthAddress = strings.Trim(p.EthAddress, " ")

	if !validate.ValidateAddress(p.EthAddress) {
		i18nresponse.Error(c, "1010001")
		return
	}

	//验证邮箱
	if r := validate.VerifyEmailFormat(p.Communication); !r {
		i18nresponse.Error(c, "1010008")
		return
	}

	created, err := participants.Create(c, p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}
	participate.Logging(created)

	i18nresponse.Success(c, "ok", created)
}

// 领取活动发送短信验证码
func ReceiveSmsSend(c *gin.Context) {

	code := c.Request.FormValue("code")
	phone := c.Request.FormValue("phone")
	languageCode := c.Request.FormValue("languageCode")
	userId := c.Request.FormValue("userId")

	if code == "" || phone == "" || languageCode == "" || userId == "" {
		i18nresponse.Error(c, "1010003")
		return
	}
	fmt.Color.Println("语言")
	fmt.Color.Println(i18nresponse.GetLang(c))
	var p participant.Participant
	p.Type = 1
	p.AreaCode = code
	p.Communication = phone
	p.Language = languageCode
	p.UserId = userId

	// 判断此账号是否已领取过 如果已领取过的那么当前账号是否与之前已绑定账号一致
	r, err := p.FindByUserId()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if r.Id != 0 {
		if r.BindAreaCode+r.Communication != code+phone {
			i18nresponse.Error(c, "1010026")
			return
		}
	}

	result, err := p.FindByCommunication()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if result.Id == 0 {
		i18nresponse.Error(c, "1010021")
		return
	}
	if result.ConvertType == 1 {
		i18nresponse.Error(c, "1010024")
		return
	}
	if result.UserId != "" && result.UserId != userId {
		i18nresponse.Error(c, "1010026")
		return
	}

	// 发送验证码
	smsResult, err := usercenter.SendCode(code, phone, languageCode)
	if err != nil || smsResult.Data.Succeed != true {
		if smsResult.Data.Message == "" {
			i18nresponse.Error(c, "1010005")
		} else {
			i18nresponse.Error(c, smsResult.Data.Message)
		}
		return
	}

	i18nresponse.Success(c, "ok", struct {
		Success bool `json:"success"`
	}{Success: true})
}

// 领取活动验证短信验证码
func ReceiveValidateCode(c *gin.Context) {

	var p participant.Participant
	var req oauth.ValidateCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	if req.ReceiveType == 2 {
		result, err := usercenter.ValidateCode(req)
		if err != nil {
			i18nresponse.Error(c, "500")
			return
		}

		if result.Success == false || result.Data.Succeed == false {
			i18nresponse.Error(c, result.Data.Message)
			return
		}
	}

	// 短信验证成功直接走领取操作
	stu, err := json.Marshal(req)
	if err = json.Unmarshal(stu, &p); err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	p.EthAddress = strings.Trim(p.EthAddress, " ")

	if !validate.ValidateAddress(p.EthAddress) {
		i18nresponse.Error(c, "1010001")
		return
	}

	convertAmount, ethAddress, err := participants.Receive(p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", struct {
		ConvertAmount int64  `json:"convert_amount"`
		EthAddress    string `json:"eth_address"`
	}{ConvertAmount: convertAmount, EthAddress: ethAddress})
}

// 领取活动发送邮箱验证码
func ReceiveSendEmailCode(c *gin.Context) {
	fmt.Color.Println("发送邮箱验证码请求参数")
	var req oauth.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}
	fmt.Color.Printf(req)
	fmt.Color.Println("请求语言")
	fmt.Color.Println(i18nresponse.GetLang(c))
	req.Email = strings.Trim(req.Email, " ")
	arr := strings.Split(req.Email, "@")
	if val := array.StrInArray(constant.Emails, arr[1]); !val {
		i18nresponse.Error(c, "_S000032")
		return
	}

	var p participant.Participant
	p.Type = 2
	p.Communication = req.Email
	p.UserId = req.UserId

	// 判断此账号是否已领取过 如果已领取过的那么当前账号是否与之前已绑定账号一致
	r, err := p.FindByUserId()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if r.Id != 0 {
		if req.Email != r.Communication {
			i18nresponse.Error(c, "1010026")
			return
		}
	}

	result, err := p.FindByCommunication()
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	if result.Id == 0 {
		i18nresponse.Error(c, "1010022")
		return
	}

	if result.ConvertType == 1 {
		i18nresponse.Error(c, "1010024")
		return
	}

	if result.UserId != "" && result.UserId != req.UserId {
		i18nresponse.Error(c, "1010026")
		return
	}

	if _, err := oauth2.SendEmailCode(req); err != nil {
		i18nresponse.Error(c, "1010002")
		return
	}

	i18nresponse.Success(c, "ok", nil)
}

// 领取验证邮箱验证码
func ReceiveValidateEmailCode(c *gin.Context) {

	var p participant.Participant
	var req oauth.ValidateEmailCode
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	req.Communication = strings.Trim(req.Communication, " ")
	// 过滤国内邮箱
	arr := strings.Split(req.Communication, "@")
	if val := array.StrInArray(constant.Emails, arr[1]); !val {
		i18nresponse.Error(c, "_S000032")
		return
	}

	if val, err := oauth2.ValidateEmailCode(req.Communication, req.Code); err != nil || !val {
		i18nresponse.Error(c, "1010004")
		return
	}

	// 邮箱验证成功直接走创建
	stu, err := json.Marshal(req)
	if err = json.Unmarshal(stu, &p); err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	p.EthAddress = strings.Trim(p.EthAddress, " ")

	if !validate.ValidateAddress(p.EthAddress) {
		i18nresponse.Error(c, "1010001")
		return
	}

	//验证邮箱
	if r := validate.VerifyEmailFormat(p.Communication); !r {
		i18nresponse.Error(c, "1010008")
		return
	}

	convertAmount, ethAddress, err := participants.Receive(p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", struct {
		ConvertAmount int64  `json:"convert_amount"`
		EthAddress    string `json:"eth_address"`
	}{ConvertAmount: convertAmount, EthAddress: ethAddress})
}

// 获取用户信息
func UserInfo(c *gin.Context) {

	var req oauth.GetUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	result, err := usercenter.UserInfo(req)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	if result.Succeed != true || result.Message != "success" {
		i18nresponse.Error(c, result.Message)
		return
	}

	i18nresponse.Success(c, "ok", result.Result)
}
