package database

import (
	"chords_app/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{}, &models.Song{}, &models.Artist{}, &models.SongArtist{}, &models.SongRequest{},
	)
}
