package test

import (
	"testing"
	"schedvault/config"
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
