package oauth

// 注册请求参数
type RegisterRequest struct {
	AreaFlag             string `json:"areaFlag"`
	AreaCode             string `json:"areaCode"`
	Phone                string `json:"phone"`
	Password             string `json:"password"`
	Email                string `json:"email"`
	Sex                  int    `json:"sex"`
	Lastname             string `json:"lastname"`
	IsSkip               int    `json:"isSkip"`
	ApplicationType      int    `json:"applicationType"`
	LanguageCode         string `json:"languageCode"`
	CountryCode          string `json:"countryCode"`
	Version              string `json:"version"`
	Ip                   string `json:"ip"`
	RequestId            string `json:"requestId"`
	RegistrationPlatform int    `json:"registrationPlatform"`
	DeviceInformation    string `json:"deviceInformation"`
	DeviceCode           string `json:"deviceCode"`
	UserFirstName        string `json:"userFirstName"`
}

// 短信验证码验证
type ValidateCodeRequest struct {
	AreaCode          string `json:"area_code"`
	PhoneNumber       string `json:"phoneNumber" binding:"required"`
	Smscode           string `json:"smscode"`
	UserId            string `json:"userId"`
	EthAddress        string `json:"eth_address" binding:"required"`
	Type              int    `json:"type" binding:"required"`
	Communication     string `json:"communication" binding:"required"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
	Event             string `json:"event"`
	ReceiveType       int    `json:"receive_type"`
}

// 邮箱验证码验证
type ValidateEmailCode struct {
	Code          string `json:"code"`
	EthAddress    string `json:"eth_address" binding:"required"`
	Type          int    `json:"type" binding:"required"`
	Communication string `json:"communication" binding:"required"`
	//RecaptchaResponse string `json:"g-recaptcha-response" binding:"required"`
	Event string `json:"event"`
}

// 快捷登录Request
type QuickLoginRequest struct {
	AreaCode        string `json:"areaCode"`
	LanguageCode    string `json:"languageCode"`
	Phone           string `json:"phone"`
	MsgCode         string `json:"msgCode"`
	ApplicationType int    `json:"applicationType"`
	Ip              string `json:"ip"`
	Version         string `json:"version"`
	EquipmentType   int    `json:"equipmentType"`
}

// 账号密码登录Request
type LoginRequest struct {
	Account         string `json:"account"`
	Password        string `json:"password"`
	LanguageCode    string `json:"languageCode"`
	CountryCode     string `json:"countryCode"`
	Ip              string `json:"ip"`
	EquipmentType   int    `json:"equipmentType"`
	ApplicationType int    `json:"applicationType"`
}

// 通过手机号找回密码Request
type ModifyPassByPhoneRequest struct {
	UserId    string `json:"userId"`
	AreaCode  string `json:"areaCode"`
	Phone     string `json:"phone"`
	Npwd      string `json:"npwd"`
	RequestId string `json:"requestId"`
}

// 通过旧密码改新密码Request
type ModifyPassByOldRequest struct {
	UserId string `json:"userId"`
	Opwd   string `json:"opwd"`
	Npwd   string `json:"npwd"`
}

// 发送邮箱验证码Request
type SendEmailCodeRequest struct {
	Email  string `json:"email"`
	UserId string `json:"userId"`
}

// 验证邮箱（验证码）
type ConfirmEmailByCodeRequest struct {
	Email           string `json:"email"`
	UserId          string `json:"userId"`
	Code            string `json:"code"`
	ApplicationType int    `json:"applicationType"`
}

// 验证邮箱验证码
type ValidateEmailCodeRequest struct {
	Email           string `json:"email"`
	UserId          string `json:"userId"`
	Code            string `json:"code"`
	ApplicationType int    `json:"applicationType"`
}

// 验证邮箱（链接）
type ConfirmEmailByLineRequest struct {
	Email           string `json:"email"`
	UserId          string `json:"userId"`
	ApplicationType int    `json:"applicationType"`
}

// 校验邮箱是否验证Request
type CheckMailboxRequest struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

//获取用户信息Request
type GetUserInfoRequest struct {
	UserId          string `json:"userId"`
	CountryCode     string `json:"countryCode"`
	ApplicationType int    `json:"applicationType"`
}

// token
type TokenResponse struct {
	Authorization string `json:"authorization"`
}

// tokenAndUserInfo
type TokenAndUserInfoResponse struct {
	UserInfo      userInfo `json:"user_info"`
	Authorization string   `json:"authorization"`
}

type userInfo struct {
	UserId               string `json:"userId"`
	Nickname             string `json:"nickname"`
	Avatar               string `json:"avatar"`
	Sex                  int    `json:"sex"`
	Areaflag             string `json:"areaflag"`
	Areacode             string `json:"areacode"`
	Phone                string `json:"phone"`
	Email                string `json:"email"`
	Shoppingaddresscount int    `json:"shoppingaddresscount"`
	Realname             string `json:"realname"`
	Isphonecomfirm       bool   `json:"isphonecomfirm"`
	Isemailcomfirm       bool   `json:"isemailcomfirm"`
}
