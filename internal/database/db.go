package database

import (
	"chords_app/internal/config"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase(config *config.DB) (*gorm.DB, error) {
	const op = "database.db.SetupDatabase"

	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, err
}
