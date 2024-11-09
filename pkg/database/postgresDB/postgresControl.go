package postgresDB

import (
	"fmt"
	"github.com/nehachuha1/mynotes-project/pkg/services/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	database *gorm.DB
}

func makeDsn(cfg *config.Config) string {
	dsn := "postgres://" + cfg.PostgresConfig.PostgresUser + ":" + cfg.PostgresConfig.PostgresPassword + "@" +
		cfg.PostgresConfig.PostgresAddress + ":" + cfg.PostgresConfig.PostgresPort + "/" +
		cfg.PostgresConfig.PostgresDB
	return dsn
}

func NewPostgresDB(cfg *config.Config) *PostgresDatabase {
	newDsn := makeDsn(cfg)
	dbConn, err := gorm.Open(postgres.Open(newDsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("can't open gorm.Open: %v", err))
	}
	return &PostgresDatabase{
		database: dbConn,
	}
}
