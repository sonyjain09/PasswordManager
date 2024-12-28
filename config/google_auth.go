package config

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"os"
	"schedvault/models"
	"google.golang.org/api/calendar/v3"
	"time"

)

var GoogleOauthConfig *oauth2.Config

func InitGoogleAuth() {
	fmt.Println("Initializing Google OAuth configuration...")

	// Attempt to open the credentials file
	credentialsFile, err := os.Open("credentials.json")
	if err != nil {
		fmt.Printf("Error opening credentials file: %v\n", err)
		return
	}
	defer credentialsFile.Close()

	// Read and parse the credentials file
	credentialsData, err := ioutil.ReadAll(credentialsFile)
	if err != nil {
		fmt.Printf("Error reading credentials file: %v\n", err)
		return
	}

	GoogleOauthConfig, err = google.ConfigFromJSON(credentialsData, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		fmt.Printf("Error parsing credentials file: %v\n", err)
		return
	}

	fmt.Println("Google OAuth configuration initialized successfully")
}


func GetAuthURL() string {
	if GoogleOauthConfig == nil {
		fmt.Println("Google OAuth configuration is not initialized")
		return ""
	}
	return GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}


func ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	if GoogleOauthConfig == nil {
		return nil, fmt.Errorf("Google OAuth configuration is not initialized")
	}

	// Exchange the code for a token
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		return nil, fmt.Errorf("Failed to exchange authorization code: %v", err)
	}

	fmt.Printf("Token successfully exchanged: %+v\n", token)
	return token, nil
}


func SaveTokenToDB(userID uint, token *oauth2.Token) error {
	// Create a new GoogleToken instance
	googleToken := models.GoogleToken{
		UserID:       userID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	// Save to database
	if err := DB.Create(&googleToken).Error; err != nil {
		fmt.Printf("Error saving token to DB: %v\n", err)
		return err
	}

	fmt.Println("Google token saved successfully!")
	return nil
}

func GetTokenFromDB(userID uint) (*oauth2.Token, error) {
	var googleToken models.GoogleToken

	// Query the database for the token associated with the user
	if err := DB.Where("user_id = ?", userID).First(&googleToken).Error; err != nil {
		return nil, err
	}

	// Convert the database token back to an oauth2.Token
	token := &oauth2.Token{
		AccessToken:  googleToken.AccessToken,
		RefreshToken: googleToken.RefreshToken,
		TokenType:    googleToken.TokenType,
		Expiry:       googleToken.Expiry,
	}

	return token, nil
}

func FetchGoogleCalendarEvents(token *oauth2.Token) ([]*calendar.Event, error) {
	// Create an authenticated HTTP client using the token
	ctx := context.Background()
	client := GoogleOauthConfig.Client(ctx, token)

	// Create a Google Calendar service
	srv, err := calendar.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Calendar client: %v", err)
	}

	// Define the time range (e.g., now to 1 month ahead)
	now := time.Now().Format(time.RFC3339)
	oneMonthAhead := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	// Fetch events from the primary calendar
	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(now).
		TimeMax(oneMonthAhead).
		OrderBy("startTime").
		MaxResults(50).
		Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve calendar events: %v", err)
	}

	return events.Items, nil
}

