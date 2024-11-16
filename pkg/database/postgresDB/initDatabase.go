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

type RelationWorkspace struct {
	Id            int64 `gorm:"primary_key"`
	OwnerUsername string
	IsPrivate     bool
	NotesID       []int64 `gorm:"column:notes_id;type:integer[]"`
}

type RelationNote struct {
	Id            int64 `gorm:"primary_key"`
	WorkspaceID   int64
	OwnerUsername string
	NoteText      string
	IsPrivate     bool
	Tags          []string `gorm:"type:text[];column:tags"`
	CreatedAt     string
	LastEditedAt  string
}

func (pgdb *PostgresDatabase) MakeMigrations() error {
	relations := []interface{}{&RelationUser{}, &RelationRegistration{},
		&RelationWorkspace{}, &RelationNote{}}
	err := pgdb.database.AutoMigrate(relations...)
	if err != nil {
		return fmt.Errorf("failed on making migrations to database: %v", err)
	}
	return nil
}
