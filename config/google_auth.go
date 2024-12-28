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
	ctx := context.Background()
	return GoogleOauthConfig.Exchange(ctx, code)
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

func ListCalendarEvents(token *oauth2.Token) {
	ctx := context.Background()
	client := GoogleOauthConfig.Client(ctx, token)

	srv, err := calendar.New(client)
	if err != nil {
		fmt.Printf("Unable to create Calendar client: %v\n", err)
		return
	}

	// Fetch upcoming events from the primary calendar
	events, err := srv.Events.List("primary").MaxResults(10).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve calendar events: %v\n", err)
		return
	}

	// Print upcoming events
	fmt.Println("Upcoming events:")
	for _, item := range events.Items {
		// Handle missing DateTime gracefully
		start := item.Start.DateTime
		if start == "" {
			start = item.Start.Date // All-day events may use Date instead
		}
		fmt.Printf("%s: %s\n", item.Summary, start)
	}
}
