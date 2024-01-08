package database

import (
	"github.com/InnoFours/skin-savvy/config"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.ConfigDB())
	if err != nil {
		return nil, err
	}

	return db, nil
}
