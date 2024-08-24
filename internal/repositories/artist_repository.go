package repositories

import (
	"chords_app/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ArtistRepository interface {
	CreateArtist(artist *models.Artist) error
	GetArtists() (*[]models.Artist, error)
	GetArtistById(artistId uint) (*models.Artist, error)
	GetArtistSongs(artistId uint) ([]models.Song, error)
}

type gormArtistRepository struct {
	db *gorm.DB
}

func NewGormArtistRepository(db *gorm.DB) ArtistRepository {
	return &gormArtistRepository{db: db}
}

func (r *gormArtistRepository) CreateArtist(artist *models.Artist) error {
	return r.db.Create(artist).Error
}

func (r *gormArtistRepository) GetArtists() (*[]models.Artist, error) {
	var artists []models.Artist

	result := r.db.Find(&artists)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &artists, nil
	}
	return &artists, result.Error
}

func (r *gormArtistRepository) GetArtistById(artistId uint) (*models.Artist, error) {
	var artist models.Artist

	result := r.db.Where("id = ?", artistId).First(&artist)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &artist, result.Error
}

func (r *gormArtistRepository) GetArtistSongs(artistId uint) ([]models.Song, error) {
	var songs []models.Song

	result := r.db.Model(&models.Song{}).
		Joins("JOIN song_artists ON song_artists.song_id = songs.id").
		Where("song_artists.artist_id = ?", artistId).
		Preload("Artists", func(db *gorm.DB) *gorm.DB {
			return db.Order("title_order")
		}).
		Find(&songs)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(songs) == 0 {
		return nil, nil
	}
	return songs, nil
}
