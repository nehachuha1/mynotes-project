package user

import (
	"context"
	"errors"
	"fmt"
	argonpass "github.com/dwin/goArgonPass"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/database/postgresDB"
	"go.uber.org/zap"
	"time"
)

var (
	ErrBadField            = errors.New("bad field in input structure")
	ErrBadUsername         = errors.New("bad field 'Username' in input structure")
	ErrWrongPasswordLength = errors.New("password's length must be equal or more than 8 symbols")
	ErrPasswordIsNotMatch  = errors.New("password from database and current request are not equal")
)

type UserRepository struct {
	abstractions.UnimplementedAuthServiceServer
	PostgresManager *postgresDB.PostgresDatabase
	logger          *zap.SugaredLogger
}

func NewUserRepository(postgresDatabase *postgresDB.PostgresDatabase, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		PostgresManager: postgresDatabase,
		logger:          logger,
	}
}

func (repo *UserRepository) RegisterUser(ctx context.Context, newRegistration *abstractions.Registration) (*abstractions.Result, error) {
	err := validateInputRegistration(newRegistration)
	if err != nil {
		repo.logger.Warnw("bad field in structure", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    400,
			Message: fmt.Sprintf("Failed to register user in RegisterUser: %v", err),
		}
		return result, err
	} else if len(newRegistration.GetPassword()) < 8 {
		repo.logger.Warnw("bad field in structure", "type", "gRPC server endpoint",
			"output", ErrWrongPasswordLength.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    400,
			Message: fmt.Sprintf("Failed to register user in RegisterUser: %v", ErrWrongPasswordLength.Error()),
		}
		return result, ErrWrongPasswordLength
	}
	hashedPassword, err := argonpass.Hash(newRegistration.GetPassword(), nil)
	if err != nil {
		repo.logger.Warnw("can't hash password", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    500,
			Message: fmt.Sprintf("Failed to register user in RegisterUser: %v", err.Error()),
		}
		return result, err
	}
	newRegistration.Password = hashedPassword
	err = repo.PostgresManager.RegisterUser(newRegistration)
	if err != nil {
		repo.logger.Warnw("failed register user in user repository", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    500,
			Message: fmt.Sprintf("Failed to register user in RegisterUser: %v", err),
		}
		return result, err
	}
	result := &abstractions.Result{
		Code:    200,
		Message: fmt.Sprintf("Successfully registered user in RegisterUser"),
	}
	repo.logger.Infow("successfully registered user", "type", "gRPC server endpoint",
		"output", "REGISTERED USER", "time", time.Now().String())
	return result, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, newUser *abstractions.User) (*abstractions.Result, error) {
	err := validateInputUser(newUser)
	if err != nil {
		repo.logger.Warnw("bad field in structure", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    400,
			Message: fmt.Sprintf("Failed to create user in CreateUser: %v", err),
		}
		return result, err
	}
	err = repo.PostgresManager.CreateUser(newUser)
	if err != nil {
		repo.logger.Warnw("failed create user in user repository", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    500,
			Message: fmt.Sprintf("Failed to create user in CreateUser: %v", err),
		}
		return result, err
	}
	result := &abstractions.Result{
		Code:    200,
		Message: fmt.Sprintf("Successfully created user in CreateUser"),
	}
	repo.logger.Infow("successfully created user", "type", "gRPC server endpoint",
		"output", "CREATED USER", "time", time.Now().String())
	return result, nil
}

// TODO: сделать проверку по паролю в этом методе
func (repo *UserRepository) AuthorizeUser(ctx context.Context, user *abstractions.Registration) (*abstractions.Registration, error) {
	err := validateInputRegistration(user)
	if err != nil {
		repo.logger.Warnw("bad field in structure", "type", "gRPC server endpoint",
			"output", ErrBadUsername.Error(), "time", time.Now().String())
		result := &abstractions.Registration{
			Id:       0,
			Username: "",
			Password: "",
		}
		return result, err
	} else if len(user.GetPassword()) < 8 {
		repo.logger.Warnw("bad field in structure", "type", "gRPC server endpoint",
			"output", ErrWrongPasswordLength.Error(), "time", time.Now().String())
		result := &abstractions.Registration{
			Id:       0,
			Username: "",
			Password: "",
		}
		return result, ErrWrongPasswordLength
	}

	registration, err := repo.PostgresManager.AuthorizeUser(user)
	if err != nil {
		repo.logger.Warnw("failed authorize user in user repository", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Registration{
			Id:       0,
			Username: "",
			Password: "",
		}
		return result, err
	}
	isPasswordMatchErr := argonpass.Verify(user.GetPassword(), registration.GetPassword())
	if isPasswordMatchErr != nil {
		repo.logger.Warnw("failed authorize user in user repository", "type", "gRPC server endpoint",
			"output", ErrPasswordIsNotMatch.Error(), "time", time.Now().String())
		return &abstractions.Registration{}, isPasswordMatchErr
	}
	repo.logger.Infow("successfully authorized user", "type", "gRPC server endpoint",
		"output", "AUTHORIZED USER", "time", time.Now().String())
	return registration, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, user *abstractions.User) (*abstractions.Result, error) {
	err := repo.PostgresManager.DeleteUser(user)
	if err != nil {
		repo.logger.Warnw("failed create user in user repository", "type", "gRPC server endpoint",
			"output", err.Error(), "time", time.Now().String())
		result := &abstractions.Result{
			Code:    500,
			Message: fmt.Sprintf("Failed to authorize user in DeleteUser: %v", err),
		}
		return result, err
	}
	result := &abstractions.Result{
		Code:    200,
		Message: fmt.Sprintf("Successfully deleted user in DeleteUser"),
	}
	repo.logger.Infow("successfully deleted user", "type", "gRPC server endpoint",
		"output", "DELETED USER", "time", time.Now().String())
	return result, nil
}

func validateInputUser(user *abstractions.User) error {
	if user.Telegram == "" || user.Username == "" || user.Email == "" || user.Initials == "" {
		return ErrBadField
	}
	return nil
}

func validateInputRegistration(reg *abstractions.Registration) error {
	if reg.Username == "" {
		return ErrBadUsername
	}
	return nil
}
