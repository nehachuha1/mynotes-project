package user

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/database/postgresDB"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func StartServer(db *postgresDB.PostgresDatabase, logger *zap.SugaredLogger) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(fmt.Errorf("can't start server on :8081"))
	}
	server := grpc.NewServer()
	abstractions.RegisterAuthServiceServer(server, NewUserRepository(db, logger))
	logger.Info("Started listening on :8081")
	server.Serve(lis)
}

func TestNewUserRepository(t *testing.T) {
	// load .env
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("can't load .env file: %v", err))
	}
	cfg := config.NewConfig()
	// create logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	pgdb := postgresDB.NewPostgresDB(cfg, sugaredLogger)
	t.Logf("Successfully iniialized Postgres: %#v", pgdb)
	go StartServer(pgdb, sugaredLogger)

	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure())
	if err != nil {
		sugaredLogger.Fatalf("can't start gRPC client")
		panic("can't start gRPC client")
	}

	defer grpcConn.Close()
	userManager := abstractions.NewAuthServiceClient(grpcConn)
	ctx := context.Background()

	result, err := userManager.RegisterUser(ctx, &abstractions.Registration{
		Username: "testUsername1",
		Password: "testPassword1",
	})
	if err != nil {
		t.Fatalf("User is not registered | Error: %v", err.Error())
	}
	t.Logf("User successfully registered | Code: %v | Message: %v\n",
		result.Code, result.Message)

	registration, err := userManager.AuthorizeUser(ctx, &abstractions.Registration{
		Username: "testUsername1",
		Password: "testPassword1",
	})
	if err != nil {
		t.Fatalf("User is not authorized | Error: %v\n", err.Error())
	}
	t.Logf("User with ID %v and username %v successfully authorized | Password: %v",
		registration.GetId(), registration.GetUsername(), registration.GetPassword())
	result, err = userManager.CreateUser(ctx, &abstractions.User{
		Username: "testUsername1",
		Email:    "test@test.com",
		Initials: "Test T.T.",
		Telegram: "@test",
	})
	if err != nil {
		t.Fatalf("Can't create user | Error: %v", err.Error())
	}
	t.Logf("Successfully created user in table 'relation_users' | Code: %v | Message: %v\n",
		result.GetCode(), result.GetMessage())
	result, err = userManager.DeleteUser(ctx, &abstractions.User{
		Username: "testUsername1",
		Email:    "test@test.com",
		Initials: "Test T.T.",
		Telegram: "@test",
	})
	if err != nil {
		t.Fatalf("Can't delete user | Error: %v", err.Error())
	}
	t.Logf("Successfully created user in table 'relation_users' | Code: %v | Message: %v\n",
		result.GetCode(), result.GetMessage())
}
