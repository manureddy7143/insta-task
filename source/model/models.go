package model

type Users struct {
	Id        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"password"`
	Username  string `json:"username" gorm:"unique"`
}
