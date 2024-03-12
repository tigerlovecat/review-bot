package app

type Response struct {
	Code    int         `json:"status_code" example:"200"`
	Msg     string      `json:"status_message" example:"成功"`
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

type Page struct {
	List   interface{} `json:"list"`
	Count  int64       `json:"count"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
}

type PageResponse struct {
	Code    int    `json:"status_code" example:"200"`
	Msg     string `json:"status_message"`
	Success bool   `json:"success" example:"true"`
	Data    Page   `json:"data"`
}

func (res *Response) ReturnOK() *Response {
	res.Code = 200
	res.Success = true
	return res
}

func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	res.Success = false
	return res
}

func (res *PageResponse) ReturnOK() *PageResponse {
	res.Code = 200
	res.Success = true
	return res
}
