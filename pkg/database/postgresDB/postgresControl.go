package postgresDB

import (
	"errors"
	"fmt"
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type PostgresDatabase struct {
	database *gorm.DB
	logger   *zap.SugaredLogger
}

func makeDsn(cfg *config.Config) string {
	dsn := "postgres://" + cfg.PostgresConfig.PostgresUser + ":" + cfg.PostgresConfig.PostgresPassword + "@" +
		cfg.PostgresConfig.PostgresAddress + ":" + cfg.PostgresConfig.PostgresPort + "/" +
		cfg.PostgresConfig.PostgresDB
	return dsn
}

func NewPostgresDB(cfg *config.Config, logger *zap.SugaredLogger) *PostgresDatabase {
	newDsn := makeDsn(cfg)
	dbConn, err := gorm.Open(postgres.Open(newDsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("can't open gorm.Open: %v", err))
	}
	return &PostgresDatabase{
		database: dbConn,
		logger:   logger,
	}
}

func (pgdb *PostgresDatabase) RegisterUser(newRegistration *abstractions.Registration) error {
	result := &abstractions.Registration{}
	pgdb.database.Table("relation_registrations").Where("username = ?", newRegistration.Username).First(result)
	if result.GetUsername() != "" {
		pgdb.logger.Warnw("failed register user in postgres control", "type", "postgres",
			"output", "user with this username is already registered", "time", time.Now().String())
		return fmt.Errorf("user with username %v already exists", newRegistration.GetUsername())
	}
	newUser := &RelationRegistration{
		Username: newRegistration.GetUsername(),
		Password: newRegistration.GetPassword(),
	}
	resultFromDB := pgdb.database.Create(newUser)
	if resultFromDB.Error != nil {
		pgdb.logger.Warnw("failed register user in postgres control", "type", "postgres",
			"output", resultFromDB.Error, "time", time.Now().String())
		return fmt.Errorf("failed register user in postgres control: %v", resultFromDB.Error)
	}
	pgdb.logger.Infow("successfully registered user", "type", "postgres",
		"output", "REGISTERED USER IN POSTGRES", "time", time.Now().String())
	return nil
}

func (pgdb *PostgresDatabase) CreateUser(user *abstractions.User) error {
	currentUser := &RelationUser{
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Initials: user.GetInitials(),
		Telegram: user.GetTelegram(),
	}
	result := pgdb.database.Create(currentUser)
	if result.Error != nil {
		pgdb.logger.Warnw("failed creating user in postgres control", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("failed creating user in postgres control: %v", result.Error)
	}
	pgdb.logger.Infow("successfully created user", "type", "postgres",
		"output", "CREATED USER IN POSTGRES", "time", time.Now().String())
	return nil
}

func (pgdb *PostgresDatabase) AuthorizeUser(user *abstractions.Registration) (*abstractions.Registration, error) {
	userToAuthorize := &RelationRegistration{
		Username: user.GetUsername(),
	}
	result := pgdb.database.Table("relation_registrations").Where("username = ?", userToAuthorize.Username).First(userToAuthorize)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		pgdb.logger.Warnw("failed on authorizing user", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return nil, fmt.Errorf("failed on authorizing user: %v", result.Error)
	}

	checkedUser := &abstractions.Registration{
		Id:       userToAuthorize.Id,
		Username: userToAuthorize.Username,
		Password: userToAuthorize.Password,
	}
	pgdb.logger.Infow("successfully got user", "type", "postgres",
		"output", "GREP USER FROM POSTGRES", "time", time.Now().String())
	return checkedUser, nil
}

func (pgdb *PostgresDatabase) DeleteUser(user *abstractions.User) error {
	userToDelete := &RelationUser{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Initials: user.GetInitials(),
		Telegram: user.GetTelegram(),
	}
	checkedUser := &RelationUser{}
	result := pgdb.database.Table("relation_users").Where("username = ?", userToDelete.Username).First(checkedUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		pgdb.logger.Warnw("can't find user in database that could be deleted", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't find user in database that could be deleted: %v", result.Error)
	}
	result = pgdb.database.Table("relation_users").Delete(RelationUser{}, checkedUser.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("can't delete user in relation_user table", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't delete user in relation_user table: %v", result.Error)
	}
	result = pgdb.database.Table("relation_registrations").Delete(&RelationRegistration{}, checkedUser.Id)
	if result.Error != nil {
		pgdb.logger.Warnw("can't delete user in relation_registration table", "type", "postgres",
			"output", result.Error, "time", time.Now().String())
		return fmt.Errorf("can't delete user in relation_registration table: %v", result.Error)
	}
	pgdb.logger.Infow("successfully deleted user", "type", "postgres",
		"output", "DELETED USER FROM POSTGRES", "time", time.Now().String())
	return nil
}
