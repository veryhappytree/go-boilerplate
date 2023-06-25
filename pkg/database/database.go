package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go-boilerplate/config"
	"strconv"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var DbConnection *sql.DB

var ErrorDbConnectionFailed = errors.New("db connection was failed")

func Setup(config config.DatabaseConfig) {
	var err error

	p, err := strconv.Atoi(config.Port)
	if err != nil {
		log.Panic().Err(err).Msg("[SQL] db connection was failed")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.Host, config.User, config.Password, config.Name, p)

	var cfg = &gorm.Config{}

	if config.Debug > 0 {
		cfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	}

	DB, err = gorm.Open(postgres.Open(dsn), cfg)

	if DB.Error != nil || err != nil {
		log.Panic().Err(ErrorDbConnectionFailed)
	}

	DbConnection, err = DB.DB()
	if err != nil {
		log.Panic().Err(ErrorDbConnectionFailed).Msg("[SQL] db connection was failed")
	}

	var ping bool
	DB.Raw("select 1").Scan(&ping)
	if !ping {
		log.Panic().Err(ErrorDbConnectionFailed).Msg("[SQL] db connection was failed")
	}

	log.Info().Msg("[SQL] connection was successfully opened to database")
}

func EnsureMigrations(migrations []*gormigrate.Migration) {
	m := gormigrate.New(DB, &gormigrate.Options{
		TableName:                 "gorm_migrations",
		IDColumnName:              "id",
		IDColumnSize:              512,
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}, migrations)

	if err := m.Migrate(); err != nil {
		log.Fatal().Msgf("[SQL] could not migrate: %v", err)
	}

	log.Info().Msg("[SQL] migration did run successfully")
}
