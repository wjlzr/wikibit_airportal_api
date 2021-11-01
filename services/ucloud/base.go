package ucloud

var (
	Action    = "SendUSMSMessage"
	ProjectId = "org-oegnvm" // 项目ID
	//SigContent = "SIG20210301A288A8" // 国内短信签名
	TemplateId = "UTB210302E20DC2"   // 模板ID
	SigContent = "SIG20210225606530" // 国际短信签名
)

type smsResponse struct {
	Action    string `json:"Action"`
	Message   string `json:"Message"`
	RetCode   int64  `json:"RetCode"`
	SessionNo string `json:"SessionNo"`
}
