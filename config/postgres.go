package config

import (
	"fmt"
	"os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"schedvault/models" 

	"schedvault/util"
)

// database instance
var DB *gorm.DB

func init() {
	util.InitEnv()
}

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// connect to database with correct arguments
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

	//  create database tables based on the struct definitions
    database.AutoMigrate(&models.User{}, &models.Availability{}, &models.Booking{})

	// set global variable
    DB = database
}
