package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"schedvault/config"
	"schedvault/models"
	"schedvault/routes"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

)

func TestConnectDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConnectDatabase panicked: %v", r)
		}
	}()

	config.ConnectDatabase()

	if config.DB == nil {
		t.Error("Database connection failed: DB is nil")
	}
}

func TestRegisterUser(t *testing.T) {
	router := routes.SetupRouter()

	body := `{
		"email": "testuser_register@example.com",
		"password": "testpassword123"
	}`

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", "testuser_register@example.com").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestLoginUser(t *testing.T) {
	defer CleanupDatabase()

	router := routes.SetupRouter()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword123"), bcrypt.DefaultCost)
	user := models.User{
		Email:    "testuser_login@example.com",
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Logf("Test user created: %+v", user)

	body := `{
		"email": "testuser_login@example.com",
		"password": "testpassword123"
	}`

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200, got %d", w.Code)
	assert.Contains(t, w.Body.String(), "token")
}


func CleanupDatabase() {
	config.DB.Exec("DELETE FROM users WHERE email LIKE 'testuser_%'")
}