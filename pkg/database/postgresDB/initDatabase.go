package postgresDB

import (
	"fmt"
)

type RelationUser struct {
	Id       int64 `gorm:"primary_key"`
	Username string
	Email    string
	Initials string
	Telegram string
}

type RelationRegistration struct {
	Id       int64 `gorm:"primary_key"`
	Username string
	Password string
}

func (pgdb *PostgresDatabase) MakeMigrations() error {
	relations := []interface{}{&RelationUser{}, &RelationRegistration{}}
	err := pgdb.database.AutoMigrate(relations...)
	if err != nil {
		return fmt.Errorf("failed on making migrations to database: %v", err)
	}
	return nil
}
