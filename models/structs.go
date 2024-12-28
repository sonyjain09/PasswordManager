package models

import (
	"time"
)

type User struct {
    ID       uint   
    Email    string 
    Password string 
}

type Vault struct {
    ID           uint 
    UserID       uint  
    Name         string 
    MasterKey    string 
}

type Event struct {
    ID        uint  
    UserID    uint 
    Title     string
    StartTime string 
    EndTime   string 
}

type PasswordRecord struct {
    ID          uint
    VaultID     uint   
    Name        string  
    Username    string 
    EncryptedPassword string 
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type RegisterInput struct {
	Email    string 
	Password string 
}

type LoginInput struct {
	Email    string 
	Password string 
}
