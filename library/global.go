package global

type Response struct {
	Code int         `json:"code" example:"200"` // 代码
	Msg  string      `json:"msg"`                // 消息
	Data interface{} `json:"data"`               // 数据集
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

type PageResponse struct {
	Code int    `json:"code" example:"200"` // 代码
	Data Page   `json:"data"`               // 数据集
	Msg  string `json:"msg"`                // 消息
}

//成功状态
func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

//错误
func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}

func (res *PageResponse) ReturnOK() *PageResponse {
	res.Code = 200
	return res
}
