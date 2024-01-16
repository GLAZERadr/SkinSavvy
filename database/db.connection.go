package database

import (
	"log"

	"github.com/InnoFours/skin-savvy/config"
	"github.com/InnoFours/skin-savvy/models/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.ConfigDB()), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to DB")
	}
	DB.AutoMigrate(&entity.User{})

	log.Println("Database is up and running🫠")
}
