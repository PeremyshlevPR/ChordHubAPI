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
	GetPopularSongsForPeriod(periodDays, limit, offset uint) (*[]SongWithViews, error)
	CreateSong(song *models.Song) error
	GetSongById(songId uint) (*models.Song, error)
	GetSongWithArtists(songId uint) (*models.Song, error)
	UpdateSong(song *models.Song) error
	DeleteSong(song *models.Song) error
	AttachAuthor(songArtist *models.SongArtist) error
	DeattachAuthor(songArtist *models.SongArtist) error
	AddSongRequest(songId uint) error
}

type gormSongRepository struct {
	db *gorm.DB
}

func NewGormSongRepository(db *gorm.DB) SongRepository {
	return &gormSongRepository{db}
}

func (r *gormSongRepository) GetPopularSongsForPeriod(periodDays, limit, offset uint) (*[]SongWithViews, error) {
	var result []SongWithViews

	subquery := r.db.
		Select("song_id, COUNT(*) as view_count").
		Table("song_requests").
		Group("song_id")

	if periodDays > 0 {
		subquery = subquery.Where("requested_at >= ?", time.Now().AddDate(0, 0, -int(periodDays)))
	}

	err := r.db.
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

func (r *gormSongRepository) CreateSong(song *models.Song) error {
	return r.db.Create(song).Error
}

func (r *gormSongRepository) GetSongById(songId uint) (*models.Song, error) {
	var song models.Song
	err := r.db.Where("id = ?", songId).First(&song).Error
	return &song, err
}

func (r *gormSongRepository) GetSongWithArtists(songId uint) (*models.Song, error) {
	var song models.Song

	err := r.db.Model(&models.Song{}).
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

func (r *gormSongRepository) UpdateSong(song *models.Song) error {
	return r.db.Save(song).Error
}

func (r *gormSongRepository) DeleteSong(song *models.Song) error {
	return r.db.Delete(song).Error
}

func (r *gormSongRepository) AttachAuthor(songArtist *models.SongArtist) error {
	return r.db.Create(songArtist).Error
}

func (r *gormSongRepository) DeattachAuthor(songArtist *models.SongArtist) error {
	return r.db.Delete(songArtist).Error
}

func (r *gormSongRepository) AddSongRequest(songId uint) error {
	songRequest := models.SongRequest{
		SongID: songId,
	}
	return r.db.Create(&songRequest).Error
}
