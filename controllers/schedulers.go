package controllers

import (
	"net/http"
	"schedvault/config"
	"schedvault/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
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
	// bind JSON request to booking input struct
	var input models.BookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// get the user_id from the context
	user_id := getUserID(c)

	// checks for overlapping bookings for the same user_id
	var existing_booking models.Booking
	if err := config.DB.Where("start_time < ? AND end_time > ?", input.EndTime, input.StartTime).
		Where("user_id = ?", user_id). 
		First(&existing_booking).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No overlapping booking found")
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for booking conflicts"})
			return
		}
	} else {
		fmt.Printf("Overlapping booking found: %+v\n", existing_booking)
		c.JSON(http.StatusConflict, gin.H{"error": "Time slot already booked"})
		return
	}

	// if no overlapping booking create new struct to add to database
	booking := models.Booking{
		UserID:       user_id,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		BookedBy:     input.BookedBy,
		BookedByEmail: input.BookedByEmail,
	}

	// save booking to database
	if err := config.DB.Create(&booking).Error; err != nil {
		fmt.Printf("Error creating booking: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// set success status
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