package usercenter

import "github.com/gin-gonic/gin"

const (
	userName               = "gsw"
	password               = "2E6CAA66096F1D4BED0FE21EE5468FE2"
	getToken               = "api/Permission/Login"                       //获取token
	sendCode               = "PersonCenter/usercenter/sendcode"           //发送短信验证码
	validateCode           = "PersonCenter/usercenter/validatecode"       //验证短信验证码
	sendEmail              = "PersonCenter/usercenter/sendemail"          //发送邮箱验证码
	validateEmailCode      = "PersonCenter/usercenter/validateemailcode"  //验证邮箱验证码
	sendCodeBit            = "Third/smsweb/sendcodebit"                   //发送短信验证码
	getCodesByRequestIdBit = "Third/smsweb/getcodesbyrequestidBit"        //校验短信验证码
	getUser                = "PersonCenter/usercenter/wikiglobal/getuser" //获取用户信息
)

var (
	Authorization   string
	c               *gin.Context
	applicationType = 61
)

type tokenRequest struct {
	UserName string
	Password string
}

// tokenRequest
type tokenResponse struct {
	Status      bool   `json:"status"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// 验证手机号response new
type ValidateUserPhoneResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"Success"`
	Data    struct {
		Result  string `json:"result"`
		Succeed bool   `json:"succeed"`
		Message string `json:"message"`
	} `json:"Data"`
}

// 发送验证码request
type sendCodeRequest struct {
	AreaCode        string `json:"areaCode"`
	Phone           string `json:"phone"`
	LanguageCode    string `json:"languageCode"`
	UserId          string `json:"userId"`
	SmsBusinessType int    `json:"smsBusinessType"`
	ApplicationType int    `json:"applicationType"`
}

// 发送验证码response
type SendCodeResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"Success"`
	Data    struct {
		Result struct {
			Requestid string `json:"requestid"`
		} `json:"Result"`
		Succeed bool   `json:"succeed"`
		Message string `json:"message"`
	} `json:"Data"`
}

// 带有用户信息通用返回的参数
type CurrencyWithUserResponse struct {
	Code    int                          `json:"code"`
	Success bool                         `json:"Success"`
	Msg     string                       `json:"msg"`
	Data    CurrencyWithUserResponseData `json:"data"`
}

type CurrencyWithUserResponseData struct {
	Succeed bool                           `json:"succeed"`
	Message string                         `json:"message"`
	Result  CurrencyWithUserResponseResult `json:"result"`
}

type CurrencyWithUserResponseResult struct {
	UserId               string `json:"userId"`
	Nickname             string `json:"nickname"`
	Avatar               string `json:"avatar"`
	Sex                  int    `json:"sex"`
	Areaflag             string `json:"areaflag"`
	Areacode             string `json:"areacode"`
	Phone                string `json:"phone"`
	RealPhone            string `json:"realphone"`
	Email                string `json:"email"`
	Shoppingaddresscount int    `json:"shoppingaddresscount"`
	Realname             string `json:"realname"`
	Isphonecomfirm       bool   `json:"isphonecomfirm"`
	Isemailcomfirm       bool   `json:"isemailcomfirm"`
}

//验证短信验证码Response old
type ValidateCodeResponseOld struct {
	RequestId string `json:"RequestId"`
	Timestamp string `json:"Timestamp"`
	Content   struct {
		Result struct {
			Succeed bool   `json:"succeed"`
			Message string `json:"message"`
		} `json:"result"`
	} `json:"Content"`
}

// 通用response new
type CurrencyResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"Success"`
	Data    struct {
		Succeed bool   `json:"succeed"`
		Message string `json:"message"`
	} `json:"Data"`
}
