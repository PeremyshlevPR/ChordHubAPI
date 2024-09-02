package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string
	Email         string
	PasswordHash  string
	Role          string
	UploadedSongs []Song `gorm:"foreignKey:UploadedBy;constraint:OnDelete:CASCADE;"`
}

type Artist struct {
	gorm.Model
	Name        string
	Description string
	ImageUrl    string
	Songs       []SongArtist `gorm:"constraint:OnDelete:CASCADE;"`
}

type Song struct {
	gorm.Model
	Title       string
	Description string
	Content     string
	Artists     []SongArtist `gorm:"constraint:OnDelete:CASCADE;"`
	UploadedBy  uint
}

type SongArtist struct {
	gorm.Model
	ArtistID   uint
	Artist     Artist `gorm:"constraint:OnDelete:CASCADE;"`
	SongID     uint
	Song       Song `gorm:"constraint:OnDelete:CASCADE;"`
	TitleOrder int
}

type SongRequest struct {
	gorm.Model
	SongID uint
}
