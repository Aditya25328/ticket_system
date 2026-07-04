package config

import (
	"log"

	"ticket-system/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("database/ticket.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database",err)
	}

	database.AutoMigrate(
		&models.User{},
		&models.Ticket{},
	)

	DB = database

	log.Println("Database Connected")
}