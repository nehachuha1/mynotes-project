package session

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"testing"
)

func TestSessionManager(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("can't load .env file: %v", err))
	}
	newConfig := config.NewConfig()
	sessionManager := NewSessionManager(newConfig)

	//test1
	user := abstractions.User{Username: "username"}
	createdSession, err := sessionManager.CreateSession(user.Username)
	if err != nil {
		t.Fatalf("Failed on creating session in test1: %v\n", err)
	}
	newToken, err := sessionManager.CreateNewToken(&user, createdSession.SessionID)
	if err != nil {
		t.Fatalf("Failed on creating new token in test1: %v", err)
	}
	t.Logf("New token: %v\n", newToken)
}
