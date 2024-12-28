package util 

import (
	"github.com/joho/godotenv"
	"path/filepath"
	"fmt"
)

// function to read from the .env file
func InitEnv() {
	env_path, _ := filepath.Abs("./.env")
	err := godotenv.Load(env_path)
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}