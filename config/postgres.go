package config

import (
	"fmt"
	"os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"github.com/joho/godotenv"
	"schedvault/models" 
)

var DB *gorm.DB

func ConnectDatabase() {
    err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }
    database.AutoMigrate(&models.User{}, &models.Vault{}, &models.Event{})
    DB = database
}
