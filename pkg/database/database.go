package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go-boilerplate/config"
	"strconv"

	"log/slog"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var DBConnection *sql.DB

var ErrDBConnectionFailed = errors.New("db connection was failed")

func Setup(cfg *config.DatabaseConfig) {
	var err error

	p, err := strconv.Atoi(cfg.Port)
	if err != nil {
		slog.Error("[SQL] parse config", "error", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Host, cfg.User, cfg.Password, cfg.Name, p)

	var gormCfg = &gorm.Config{}

	if cfg.Debug > 0 {
		gormCfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormCfg)

	if DB.Error != nil || err != nil {
		slog.Error("[SQL] ErrorDBConnectionFailed", "error", ErrDBConnectionFailed)
		panic(ErrDBConnectionFailed)
	}

	DBConnection, err = DB.DB()
	if err != nil {
		slog.Error("[SQL] ErrorDBConnectionFailed", "error", ErrDBConnectionFailed)
		panic(ErrDBConnectionFailed)
	}

	var ping bool
	DB.Raw("select 1").Scan(&ping)
	if !ping {
		slog.Error("[SQL] db connection was failed", "error", ErrDBConnectionFailed)
		panic(ErrDBConnectionFailed)
	}

	slog.Info("[SQL]", "message", "connection was successfully opened to database")
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
		slog.Error("[SQL] could not migrate", "err", err)
		panic(err)
	}

	slog.Info("[SQL]", "message", "migration did run successfully")
}
