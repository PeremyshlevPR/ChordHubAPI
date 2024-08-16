package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string
	Email         string
	PasswordHash  string
	Role          string
	UploadedSongs []Song `gorm:"foreignKey:UploadedBy"`
}

type Artist struct {
	gorm.Model
	Name        string
	Description string
	ImageUrl    string
	Songs       []SongArtist
}

type Song struct {
	gorm.Model
	Title       string
	Description string
	Content     string
	Artists     []SongArtist
	UploadedBy  uint
}

type SongArtist struct {
	gorm.Model
	ArtistID   uint
	Artist     Artist
	SongID     uint
	Song       Song
	TitleOrder int
}
