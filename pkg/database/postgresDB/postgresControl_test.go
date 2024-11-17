package postgresDB

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nehachuha1/mynotes-project/internal/abstractions"
	"github.com/nehachuha1/mynotes-project/internal/config"
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
	err := pgDB.MakeMigrations()
	if err != nil {
		t.Fatalf("error by making migrations")
	}
	t.Logf("Successfully made migrations")

	// basic tests
	newRegistration := &abstractions.Registration{
		Username: "testUsername",
		Password: "testPassword",
	}
	err = pgDB.RegisterUser(newRegistration)
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

	// other tests
	err = pgDB.RegisterUser(newRegistration)
	if err != nil {
		t.Fatalf("failed on RegisterUser: %v", err)
	}
	err = pgDB.RegisterUser(newRegistration)
	if errors.Is(err, nil) {
		t.Fatalf("registered user in database that should't be registered")
	}
	t.Log("RegisterUser tests passed")

	err = pgDB.CreateUser(newUser)
	if err != nil {
		t.Fatalf("failed on CreateUser: %v", err)
	}
	err = pgDB.CreateUser(newUser)
	if errors.Is(err, nil) {
		t.Log("created user that shouldn't be created")
	}
	t.Log("CreateUser tests passed")

	registeredUser, err = pgDB.AuthorizeUser(newRegistration)
	if err != nil {
		t.Fatalf("failed on AuthorizeUser: %v", err)
	}
	t.Logf("userToAuthorize: %#v | Authorized user: %#v", userToAuthorize, registeredUser)
	registeredUser, err = pgDB.AuthorizeUser(&abstractions.Registration{
		Username: "UserThatIsNotExists",
		Password: "",
	})
	if errors.Is(err, nil) {
		t.Fatalf("authorized not existing user")
	}
	t.Log("Authorize tests passed")

	err = pgDB.DeleteUser(newUser)
	if err != nil {
		t.Fatalf("can't delete user with error: %v", err)
	}
	err = pgDB.DeleteUser(&abstractions.User{
		Username: "UserThatIsNotExists",
		Email:    "test@test.com",
		Initials: "Test T.T.",
		Telegram: "@test",
	})
	if errors.Is(err, nil) {
		t.Fatalf("deleted user that is not exists")
	}
}
