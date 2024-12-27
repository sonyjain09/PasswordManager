package config

import (
	"fmt"
	"os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"github.com/joho/godotenv"
	"schedvault/models" 
	"path/filepath"
)

var DB *gorm.DB

func ConnectDatabase() {
	env_path, _ := filepath.Abs("../.env")
	err := godotenv.Load(env_path)
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
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
