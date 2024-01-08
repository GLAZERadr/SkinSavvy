package auth 

import (
	"firebase.google.com/go/v4/auth"
	"github.com/jinzhu/gorm"
)

type AuthService struct {
	DB			*gorm.DB
	FireAuth	*auth.Client
}