package main

import (
	"fmt"
	"schedvault/config" 
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	config.ConnectDatabase()
	fmt.Println("Database connected successfully!")
}
