package routes

import (
	"schedvault/controllers"
	"schedvault/middleware"
	"schedvault/config"
	"net/http"

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
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
			return
		}
	
		token, err := config.ExchangeCodeForToken(code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
			return
		}
	
		// Get the logged-in user (assuming you have user info in the session or context)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	
		// Save the token to the database
		err = config.SaveTokenToDB(userID.(uint), token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
			return
		}
	
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