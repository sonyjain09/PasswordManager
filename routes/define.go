package routes

import (
	"schedvault/controllers"
	"schedvault/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// create new Gin router with default middleware
	r := gin.Default()

	// When POST endpoint at /register is hit, call the RegisterUser function
	r.POST("/register", controllers.RegisterUser)

	// When POST endpoint at /login is hit, call the LoginUser function
	r.POST("/login", controllers.LoginUser)

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