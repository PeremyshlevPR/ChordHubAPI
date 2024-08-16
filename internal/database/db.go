package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	const op = "database.db.SetupDatabase"
	const path = "../../data/chords_app.db"

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, err
}
