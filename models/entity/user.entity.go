package entity

import "time"

type User struct {
	ID			string		`json:"id"`
	Fullname	string 		`json:"fullname"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
	Photo		string		`json:"photo"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}