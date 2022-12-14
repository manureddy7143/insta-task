package dto

// Schema for registering
type RequestRegister struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Schema for login

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//ErrorDTO dto
type ErrorDTO struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// Schema for Respose Profile
type RespProfile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
