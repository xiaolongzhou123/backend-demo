package typing

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResp(code int, mess string, data interface{}) *Resp {
	return &Resp{
		Code:    code,
		Message: mess,
		Data:    data,
	}
}
