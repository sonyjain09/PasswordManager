package models

import (
	"time"
)

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Email    string `gorm:"not null;unique" json:"email"`
    Password string `gorm:"not null" json:"password"`
}

type Availability struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    DayOfWeek string    `gorm:"not null" json:"day_of_week"`
    StartTime string    `gorm:"not null" json:"start_time"`
    EndTime   string    `gorm:"not null" json:"end_time"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Booking struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    UserID       uint      `gorm:"not null" json:"user_id"`
    StartTime    time.Time `gorm:"not null" json:"start_time"`
    EndTime      time.Time `gorm:"not null" json:"end_time"`
    BookedBy     string    `gorm:"not null" json:"booked_by"`
    BookedByEmail string   `gorm:"not null" json:"booked_by_email"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BookingInput struct {
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	BookedBy     string    `json:"booked_by" binding:"required"`
	BookedByEmail string   `json:"booked_by_email" binding:"required"`
}
