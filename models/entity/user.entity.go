package entity

type User struct {
	ID			string  `json:"id"`
	Fullname	string 	`json:"fullname"`
	Email		string	`json:"email"`
	Age			int     `json:"age"`
	Password	string	`json:"password"`
}