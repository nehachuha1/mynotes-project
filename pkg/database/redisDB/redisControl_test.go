package redisDB

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"testing"
)

func TestNewRedisDatabase(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("can't load .env file: %v", err))
	}
	newConfig := config.NewConfig()
	newRedisDB := NewRedisDatabase(newConfig)

	newConn := newRedisDB.RedisConnection.Get()
	_, err := newConn.Do("SET", "test", "test")
	if err != nil {
		t.Fatalf("Can't get new connection: %v", err)
	}
	_, _ = newConn.Do("DEL", "test")

	currentSession := abstractions.Session{SessionID: "123123", Username: "username"}
	_, err = newRedisDB.CreateSession(&currentSession)
	if err != nil {
		t.Fatalf("Error by creating new session: %v", err)
	}
	_, err = newRedisDB.CheckSession(&currentSession)
	if err != nil {
		t.Fatalf("Error by checking session: %v", err)
	}
	err = newRedisDB.DeleteSession(&currentSession)
	if err != nil {
		t.Fatalf("Error by deleting session: %v", err)
	}

	sessionsToCreate := []*abstractions.Session{
		{Username: "username1"},
		{Username: "username2"},
		{Username: "username3"},
	}
	sessionsToCheck := []*abstractions.Session{
		{Username: "username4", SessionID: RandStringRunes(32)},
		sessionsToCreate[0],
		sessionsToCreate[1],
		sessionsToCreate[2],
	}
	sessionsToDelete := []*abstractions.Session{
		sessionsToCheck[0],
		sessionsToCheck[1],
		sessionsToCheck[2],
		sessionsToCheck[3],
		{Username: "123123", SessionID: "kkkk"},
	}

	for i, session := range sessionsToCreate {
		newSession, err := newRedisDB.CreateSession(session)
		if err != nil {
			t.Logf("Failed to create new session: %v", err)
		}
		sessionsToCreate[i] = newSession
	}

	for _, session := range sessionsToCheck {
		_, err := newRedisDB.CheckSession(session)
		if err != nil {
			t.Logf("Failed to check session with ID %v: %v", session.SessionID, err)
		}
	}

	for _, session := range sessionsToDelete {
		err = newRedisDB.DeleteSession(session)
		if err != nil {
			t.Logf("Failed to delete session with ID %v: %v", session.SessionID, err)
		}
	}
}
