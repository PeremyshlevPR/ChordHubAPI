package repositories

import (
	"chords_app/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SongWithViews struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	ViewCount uint
	Artists   []models.SongArtist `gorm:"foreignKey:SongID"`
}

type SongRepository interface {
	GetPopularSongsForPeriod(db *gorm.DB, periodDays, limit, offset uint) (*[]SongWithViews, error)
	CreateSong(db *gorm.DB, song *models.Song) error
	GetSongById(db *gorm.DB, songId uint) (*models.Song, error)
	GetSongWithArtists(db *gorm.DB, songId uint) (*models.Song, error)
	UpdateSong(db *gorm.DB, song *models.Song) error
	DeleteSong(db *gorm.DB, song *models.Song) error
	AttachAuthor(db *gorm.DB, songArtist *models.SongArtist) error
	DeattachAuthor(db *gorm.DB, songArtist *models.SongArtist) error
	AddSongRequest(db *gorm.DB, songId uint) error
}

type gormSongRepository struct{}

func NewGormSongRepository() SongRepository {
	return &gormSongRepository{}
}

func (r *gormSongRepository) GetPopularSongsForPeriod(db *gorm.DB, periodDays, limit, offset uint) (*[]SongWithViews, error) {
	var result []SongWithViews

	subquery := db.
		Select("song_id, COUNT(*) as view_count").
		Table("song_requests").
		Group("song_id")

	if periodDays > 0 {
		subquery = subquery.Where("requested_at >= ?", time.Now().AddDate(0, 0, -int(periodDays)))
	}

	err := db.
		Select("songs.*, COALESCE(subquery.view_count, 0) as view_count").
		Table("songs").
		Joins("LEFT JOIN (?) as subquery ON songs.id = subquery.song_id", subquery).
		Preload("Artists", func(db *gorm.DB) *gorm.DB {
			return db.Order("title_order")
		}).
		Order("view_count DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&result).Error

	return &result, err
}

func (r *gormSongRepository) CreateSong(db *gorm.DB, song *models.Song) error {
	return db.Create(song).Error
}

func (r *gormSongRepository) GetSongById(db *gorm.DB, songId uint) (*models.Song, error) {
	var song models.Song
	err := db.Where("id = ?", songId).First(&song).Error
	return &song, err
}

func (r *gormSongRepository) GetSongWithArtists(db *gorm.DB, songId uint) (*models.Song, error) {
	var song models.Song

	err := db.Model(&models.Song{}).
		Joins("JOIN song_artists ON song_artists.song_id = songs.id").
		Where("song_artists.song_id = ?", songId).
		Preload("Artists", func(db *gorm.DB) *gorm.DB {
			return db.Order("title_order")
		}).
		First(&song).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("song not found")
	}

	return &song, err
}

func (r *gormSongRepository) UpdateSong(db *gorm.DB, song *models.Song) error {
	return db.Save(song).Error
}

func (r *gormSongRepository) DeleteSong(db *gorm.DB, song *models.Song) error {
	return db.Delete(song).Error
}

func (r *gormSongRepository) AttachAuthor(db *gorm.DB, songArtist *models.SongArtist) error {
	return db.Create(songArtist).Error
}

func (r *gormSongRepository) DeattachAuthor(db *gorm.DB, songArtist *models.SongArtist) error {
	return db.Delete(songArtist).Error
}

func (r *gormSongRepository) AddSongRequest(db *gorm.DB, songId uint) error {
	songRequest := models.SongRequest{
		SongID: songId,
	}
	return db.Create(&songRequest).Error
}
