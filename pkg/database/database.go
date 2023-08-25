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
var DBConnection *sql.DB

var ErrorDBConnectionFailed = errors.New("db connection was failed")

func Setup(cfg *config.DatabaseConfig) {
	var err error

	p, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Panic().Err(err).Msg("[SQL] db connection was failed")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Host, cfg.User, cfg.Password, cfg.Name, p)

	var gormCfg = &gorm.Config{}

	if cfg.Debug > 0 {
		gormCfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormCfg)

	if DB.Error != nil || err != nil {
		log.Panic().Err(ErrorDBConnectionFailed)
	}

	DBConnection, err = DB.DB()
	if err != nil {
		log.Panic().Err(ErrorDBConnectionFailed).Msg("[SQL] db connection was failed")
	}

	var ping bool
	DB.Raw("select 1").Scan(&ping)
	if !ping {
		log.Panic().Err(ErrorDBConnectionFailed).Msg("[SQL] db connection was failed")
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
