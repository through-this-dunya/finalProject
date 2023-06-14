package database

import (
	"log"

	"github.com/through-this-dunya/finalProject/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	database.AutoMigrate(&model.User{})

	return Handler{database}
}