package wikibit

const (
	recharge = "Wikibit/app/recharge" // 活动充币
)

// 充币request
type RechargeRequest struct {
	AppId    string `json:"appId"`
	UserId   string `json:"userId"`
	Money    string `json:"money"`
	Currency string `json:"currency"`
	Secret   string `json:"secret"`
	Sign     string `json:"sign"`
	Nonce    string `json:"nonce"`
}

type RechargeResponse struct {
	Code    int64                `json:"code"`
	Msg     string               `json:"msg"`
	Success bool                 `json:"Success"`
	Data    RechargeResponseData `json:"Data"`
}

type RechargeResponseData struct {
	Result  bool   `json:"result"`
	Succeed bool   `json:"succeed"`
	Message string `json:"message"`
}
