package main

import (
	"fmt"
	"schedvault/config"
	"schedvault/routes"
)

func main() {
	// Recover from any panic during initialization
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Connect to the database
	config.ConnectDatabase()
	fmt.Println("Database connected successfully!")

	// Set up the router with all routes
	router := routes.SetupRouter()

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
