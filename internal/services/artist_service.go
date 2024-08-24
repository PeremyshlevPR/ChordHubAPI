package services

import (
	"chords_app/internal/models"
	"chords_app/internal/repositories"
)

type ArtistService interface {
	CreateArtist(name, description, imageUrl string) (*models.Artist, error)
}

type artistService struct {
	repo repositories.ArtistRepository
}

func NewArtistRepository(repo repositories.ArtistRepository) ArtistService {
	return &artistService{repo}
}

func (s *artistService) CreateArtist(name, description, imageUrl string) (*models.Artist, error) {
	artist := &models.Artist{
		Name:        name,
		Description: description,
		ImageUrl:    imageUrl,
	}

	err := s.repo.CreateArtist(artist)
	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (s *artistService) GetArtists() (*[]models.Artist, error) {
	return s.repo.GetArtists()
}

func (s *artistService) GetArtistInformation(artistId uint) (*models.Artist, *[]models.Song, error) {
	artist, err := s.repo.GetArtistById(artistId)
	var empty_songs *[]models.Song

	if err != nil {
		return nil, empty_songs, err
	}

	songs, err := s.repo.GetArtistSongs(artistId)
	if err != nil {
		return nil, empty_songs, err
	}

	if songs == nil {
		songs = empty_songs
	}

	return artist, songs, nil
}
