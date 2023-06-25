package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Migrations = []*gormigrate.Migration{
	{
		ID: "test_migration",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				select 1;
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`
				select 1;
			`).Error
		},
	},
}
