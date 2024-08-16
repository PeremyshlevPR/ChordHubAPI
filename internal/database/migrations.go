package database

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&User{}, &Song{}, &Artist{}, &SongArtist{},
	)
}
