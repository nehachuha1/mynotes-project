package postgresDB

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"go.uber.org/zap"
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("can't load .env file: %v", err))
	}
	cfg := config.NewConfig()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	newDSN := makeDsn(cfg)
	if newDSN == "" {
		t.Fatalf("Can't make dsn from config")
	}
	pgDB := NewPostgresDB(cfg, sugaredLogger)
	t.Logf("Database initialized: %v", pgDB)

	newRegistration := &abstractions.Registration{
		Username: "testUsername",
		Password: "testPassword",
	}
	err := pgDB.RegisterUser(newRegistration)
	if err != nil {
		t.Fatalf("failed on RegisterUser: %v", err)
	}
	newUser := &abstractions.User{
		Username: "testUsername",
		Email:    "test@test.com",
		Initials: "test",
		Telegram: "@test",
	}
	err = pgDB.CreateUser(newUser)
	if err != nil {
		t.Fatalf("failed on CreateUser: %v", err)
	}

	userToAuthorize := &abstractions.Registration{
		Username: "testUsername",
		Password: "",
	}
	registeredUser, err := pgDB.AuthorizeUser(userToAuthorize)
	if err != nil {
		t.Fatalf("failed on AuthorizeUser: %v", err)
	}
	t.Logf("userToAuthorize: %#v | Authorized user: %#v", userToAuthorize, registeredUser)
	err = pgDB.DeleteUser(newUser)
	if err != nil {
		t.Fatalf("failed on DeleteUser: %v", err)
	}
	err = pgDB.DeleteUser(newUser)
	if err != nil {
		t.Logf("can't delete user with error: %v", err)
	}
	t.Logf("successfully passed all tests")
}
