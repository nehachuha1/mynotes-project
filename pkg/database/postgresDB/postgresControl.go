package postgresDB

import (
	"github.com/nehachuha1/mynotes-project/pkg/abstractions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openConnection() {
	db, err := gorm.Open(postgres.Open("dsn"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&abstractions.User{}, &abstractions.Session{})
}
