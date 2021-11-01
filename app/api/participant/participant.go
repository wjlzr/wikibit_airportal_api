package participant

import (
	"github.com/k0kubun/pp"
	"strings"
	"wiki_bit/app/model/participant"
	"wiki_bit/app/service/participate"
	"wiki_bit/library/convert/xint64"
	"wiki_bit/library/fmt"
	"wiki_bit/library/i18nresponse"
	"wiki_bit/library/response"
	"wiki_bit/library/validate"

	"github.com/gin-gonic/gin"
)

// Verification 验证以太坊地址
func Verification(c *gin.Context) {

	ethAddress := c.Request.FormValue("eth_address")

	if ethAddress == "" || !validate.ValidateAddress(strings.Trim(ethAddress, " ")) {
		i18nresponse.Error(c, "1010001")
		return
	}
	// 查询此地址是否认购过
	result, err := participate.GetInfo(ethAddress)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}
	state := 1 // 1 未认购过 2 已认购过
	if result.Id != 0 {
		state = 2
	}
	i18nresponse.Success(c, "ok", struct {
		State int `json:"state"`
	}{State: state})
}

// Create 领取糖果
func Create(c *gin.Context) {

	var p participant.Participant
	if err := c.ShouldBindJSON(&p); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}
	p.EthAddress = strings.Trim(p.EthAddress, " ")

	if !validate.ValidateAddress(p.EthAddress) {
		i18nresponse.Error(c, "1010001")
		return
	}

	if p.Communication == "" {
		i18nresponse.Error(c, "1010007")
		return
	}

	//验证邮箱
	if isEmail := strings.Contains(p.Communication, "@"); isEmail || p.Type == 2 {
		if r := validate.VerifyEmailFormat(p.Communication); !r {
			i18nresponse.Error(c, "1010008")
			return
		}
	}
	p.Ip = c.ClientIP()
	created, err := participate.Create(c, p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	participate.Logging(created)

	i18nresponse.Success(c, "ok", created)
}

// UpdateCommunication 更新通信信息 (暂时不用)
func UpdateCommunication(c *gin.Context) {

	var p participant.Participant
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Error(c, -1, "1010003")
		return
	}

	if p.Communication == "" {
		response.Error(c, -1, "缺少通信信息")
		return
	}

	//判断是否是邮箱
	if isEmail := strings.Contains(p.Communication, "@"); isEmail {
		if r := validate.VerifyEmailFormat(p.Communication); !r {
			response.Error(c, -1, "邮箱格式有误")
			return
		}
		p.Type = 2
	}

	result, err := participate.UpdateCommunication(p)
	if err != nil {
		response.Error(c, -1, "修改失败")
		return
	}

	response.Success(c, result, "ok")
}

// GetInfo 获取认购信息
func GetInfo(c *gin.Context) {

	ethAddress := c.Request.FormValue("eth_address")
	if ethAddress == "" {
		i18nresponse.Error(c, "1010001")
		return
	}

	result, err := participate.GetInfo(ethAddress)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}
	pp.Println("详情接口返回数据")
	pp.Println(result)
	i18nresponse.Success(c, "ok", result)
}

// Login 登录
func Login(c *gin.Context) {

	var l participant.LoginRequest

	if err := c.ShouldBindJSON(&l); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	result, err := participate.Login(l)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", result)
}

// GetHomeInfo 获取首页数据
func GetHomeInfo(c *gin.Context) {

	result, err := participate.GetHomeInfo()
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", result)
}

// Query 给市场部查询数据统计使用
func Query(c *gin.Context) {

	type1 := c.Request.FormValue("type")
	startTime := c.Request.FormValue("start_time")
	endTime := c.Request.FormValue("end_time")
	areaCode := c.Request.FormValue("area_code")
	channelId := c.Request.FormValue("channel_id")

	if xint64.StrToInt64(type1) == 0 || xint64.StrToInt64(startTime) == 0 || xint64.StrToInt64(endTime) == 0 {
		i18nresponse.Error(c, "1010003")
		return
	}

	result, err := participate.Query(xint64.StrToInt64(type1), xint64.StrToInt64(startTime), xint64.StrToInt64(endTime), xint64.StrToInt64(channelId), areaCode)
	if err != nil {
		i18nresponse.Error(c, "500")
		return
	}

	i18nresponse.Success(c, "ok", result)
}

// ReceiveType 领取类型
func ReceiveType(c *gin.Context) {

	var p participant.Participant
	if err := c.ShouldBindJSON(&p); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	result, err := participate.ReceiveType(p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}
	fmt.Color.Println("ReceiveType接口返回数据")
	fmt.Color.Printf(result)
	i18nresponse.Success(c, "ok", result)
}

// Recharge 领取
func Recharge(c *gin.Context) {

	var p participant.RechargeRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		i18nresponse.Error(c, "1010003")
		return
	}

	result, err := participate.Recharge(p)
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", struct {
		Received int64 `json:"received"`
	}{Received: result})
}

// Recharge 领取
func Clear(c *gin.Context) {

	userId := c.Request.FormValue("user_id")
	if userId == "" {
		i18nresponse.Error(c, "1010003")
		return
	}

	//var p participant.Participant
	result, err := participant.Participant{UserId: userId}.FindByUserId()
	if err != nil {
		i18nresponse.Error(c, err.Error())
		return
	}

	i18nresponse.Success(c, "ok", struct {
		IsCheat int64 `json:"is_cheat"`
	}{IsCheat: result.IsCheat})
}
