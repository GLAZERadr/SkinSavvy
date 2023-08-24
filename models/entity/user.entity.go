package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID			primitive.ObjectID  `gorm:"column:id;primaryKey" json:"userID" bson:"_id,omitempty"`
	Name		string 				`json:"name"`
	Email		string				`json:"email"`
	Password	string				`json:"password"`
	Handphone	int					`json:"phone_no"`
	CreatedAt	time.Time			`json:"createdAt"`
	UpdatedAt	time.Time			`json:"updatedAt"`
}