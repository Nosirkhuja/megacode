package database

import (
	"MegaCode/internal/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	connection, err := gorm.Open(mysql.Open("root:rootroot@/users"), &gorm.Config{})

	if err != nil {
		return err
	}

	DB = connection

	err = connection.AutoMigrate(&model.User{})
	return err
}
