package request 

type UserRegister struct {
	ID			string `json:"id"`
	Email		string `json:"email"`
	Fullname	string `json:"fullname"`
	Age			int    `json:"age"`
	Password	string `json:"password"`
}