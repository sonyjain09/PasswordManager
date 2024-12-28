package controllers

import (
	"net/http"
	"schedvault/config"
	"schedvault/models"

	"github.com/gin-gonic/gin"
	"fmt"
	"golang.org/x/oauth2"
	"time"
)

func DefineAvailability(c *gin.Context) {
	// bind JSON request to Availability struct
    var availability models.Availability
    if err := c.ShouldBindJSON(&availability); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// get the user_id from the context and set it in the new availability
    user_id := getUserID(c)
    availability.UserID = user_id

	// add the availability to the database
    if err := config.DB.Create(&availability).Error; err != nil {
        fmt.Printf("Error saving availability to DB: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save availability"})
        return
    }

	// set success status
    c.JSON(http.StatusOK, gin.H{"message": "Availability defined successfully"})
}


func GetAvailability(c *gin.Context) {
	// get the user_id from the context
	user_id := getUserID(c)

	// get all availabilities associated with the user_id
    var availabilities []models.Availability
    if err := config.DB.Where("user_id = ?", user_id).Find(&availabilities).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch availability"})
        return
    }

	// set success status
    c.JSON(http.StatusOK, gin.H{"availability": availabilities})
}


func BookSlot(c *gin.Context) {
	var input models.BookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check for overlapping slots in the database
	var existingBooking models.Booking
	if err := config.DB.Where("start_time < ? AND end_time > ?", input.EndTime, input.StartTime).
		Where("user_id = ?", userID).
		First(&existingBooking).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Time slot already booked in the system"})
		return
	}

	// Retrieve the user's token from the database
	var token models.GoogleToken
	if err := config.DB.Where("user_id = ?", userID).First(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Google token"})
		return
	}

	// Fetch Google Calendar events
	oauthToken := &oauth2.Token{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
	}
	googleEvents, err := config.FetchGoogleCalendarEvents(oauthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Google Calendar events"})
		return
	}

	// Check for overlaps with Google Calendar events
	for _, event := range googleEvents {
		eventStart, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		eventEnd, _ := time.Parse(time.RFC3339, event.End.DateTime)

		if eventStart.Before(input.EndTime) && eventEnd.After(input.StartTime) {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Time slot conflicts with Google Calendar event: %s", event.Summary)})
			return
		}
	}

	// Create the booking
	booking := models.Booking{
		UserID:       userID.(uint),
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		BookedBy:     input.BookedBy,
		BookedByEmail: input.BookedByEmail,
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
}


func GetBookings(c *gin.Context) {

	// get the user_id from the context
	user_id := getUserID(c)

	// get all bookings associated with the user_id
	var bookings []models.Booking
	if err := config.DB.Where("user_id = ?", user_id).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// set success status
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func getUserID(c *gin.Context) uint {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0
	} else {
		id, _ := user_id.(uint)
		return id
	}
}