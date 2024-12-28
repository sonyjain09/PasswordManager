package controllers

import (
	"net/http"
	"schedvault/config"
	"schedvault/models"
	"schedvault/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"os"
	"time"
)

func init() {
	util.InitEnv()
}

var jwt_secret = os.Getenv("JWT_SECRET")

func RegisterUser(c *gin.Context) {
	// parse input
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the plaintext password
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// save the user to the database
	user := models.User{Email: input.Email, Password: string(hashed_password)}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// set success status
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}


func LoginUser(c *gin.Context) {
	// parse input
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// checks database for user
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		fmt.Printf("User not found: %s\n", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		fmt.Printf("Password mismatch for email: %s\n", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate a jwt token (like cookies)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})	
	
	tokenString, err := token.SignedString([]byte(jwt_secret))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
