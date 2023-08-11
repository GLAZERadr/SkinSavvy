package request 

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegister struct {
	Name		string 		`json:"name"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
	CreatedAt	time.Time	`json:"created"`
	UpdatedAt	time.Time	`json:"updated"`
}

type UserLogin struct {
	ID			primitive.ObjectID  `json:"userID" bson:"_id,omitempty"`
	Email		string				`json:"email"`
	Password	string				`json:"password"`
}
