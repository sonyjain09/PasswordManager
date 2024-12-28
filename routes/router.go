package routes

import (
	"schedvault/controllers"
	"schedvault/middleware"
	"schedvault/config"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"

)

func SetupRouter() *gin.Engine {
	// create new Gin router with default middleware
	r := gin.Default()

	// root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the SchedVault API!"})
	})

	// When POST endpoint at /register is hit, call the RegisterUser function
	r.POST("/register", controllers.RegisterUser)

	// When POST endpoint at /login is hit, call the LoginUser function
	r.POST("/login", controllers.LoginUser)

	r.GET("/oauth2login", func(c *gin.Context) {
		authURL := config.GetAuthURL()
		c.Redirect(http.StatusFound, authURL)
	})

	r.GET("/oauth2callback", func(c *gin.Context) {
		// Retrieve the authorization code from the query string
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
			return
		}
	
		// Exchange the code for a token
		token, err := config.ExchangeCodeForToken(code)
		if err != nil {
			fmt.Printf("Error exchanging token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	
		// Debugging token
		fmt.Printf("Received Token: %+v\n", token)
	
		// Optionally associate with the logged-in user
		var userID uint
		if userIDRaw, exists := c.Get("user_id"); exists {
			userID = userIDRaw.(uint) // Safely cast to uint
			fmt.Printf("Associating token with user_id: %d\n", userID)
	
			// Save token in the database for the user
			err = config.SaveTokenToDB(userID, token)
			if err != nil {
				fmt.Printf("Error saving token to DB: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
				return
			}
		} else {
			fmt.Println("No user_id found in context; skipping association")
		}
	
		// Respond to the client
		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful", "token": token})
	})	
	

	// Create a route group for all endpoints under the /protected prefix.
	protected := r.Group("/protected")
	// only authenticated requests with a valid JWT token can access these endpoints.
	protected.Use(middleware.AuthMiddleware())
	{	
		protected.POST("/availability", controllers.DefineAvailability) // Define availability
		protected.GET("/availability", controllers.GetAvailability)     // Get availability
		protected.POST("/book", controllers.BookSlot) // Book a time slot
		protected.GET("/bookings", controllers.GetBookings) // Get booked slots

		// Defines a GET endpoint at /protected/profile
		protected.GET("/profile", func(c *gin.Context) {
			// send an OK response
			c.JSON(200, gin.H{"message": "You are authorized"})
		})
	}

	return r
}