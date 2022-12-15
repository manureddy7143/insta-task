package dto

// ErrorDTO dto
type ErrorDTO struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type Response struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}
