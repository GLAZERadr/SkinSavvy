package request 

type UserRegister struct {
	ID			string `json:"id"`
	Email		string `json:"email"`
	Fullname	string `json:"fullname"`
	Password	string `json:"password"`
}