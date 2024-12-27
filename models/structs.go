package models

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