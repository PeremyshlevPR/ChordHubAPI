package services

import (
	"chords_app/internal/models"
	"chords_app/internal/repositories"

	"errors"
)

type ArtistService interface {
	CreateArtist(name, description, imageUrl string) (*models.Artist, error)
	UpdateArtist(artistId uint, name, description, imageUrl string) (*models.Artist, error)
	DeleteArtist(artistId uint) error
	GetArtists() (*[]models.Artist, error)
	GetArtistInformation(artistId uint) (*models.Artist, *[]models.Song, error)
}

type artistService struct {
	repo repositories.ArtistRepository
}

func NewArtistService(repo repositories.ArtistRepository) ArtistService {
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

func (s *artistService) UpdateArtist(artistId uint, name, description, imageUrl string) (*models.Artist, error) {
	artist, err := s.repo.GetArtistById(artistId)
	if err != nil {
		return nil, err
	}
	if artist == nil {
		return nil, errors.New("artist not found")
	}

	if name != "" {
		artist.Name = name
	}
	if description != "" {
		artist.Description = description
	}
	if imageUrl != "" {
		artist.ImageUrl = imageUrl
	}

	err = s.repo.UpdateArtist(artist)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *artistService) DeleteArtist(artistId uint) error {
	artist, err := s.repo.GetArtistById(artistId)
	if err != nil {
		return err
	}
	if artist == nil {
		return errors.New("artist not found")
	}

	return s.repo.DeleteArtist(artist)
}

func (s *artistService) GetArtistInformation(artistId uint) (*models.Artist, *[]models.Song, error) {
	artist, err := s.repo.GetArtistById(artistId)
	if err != nil {
		return nil, nil, err
	}
	if artist == nil {
		return nil, nil, errors.New("artist not found")
	}

	var empty_songs *[]models.Song

	songs, err := s.repo.GetArtistSongs(artist.ID)
	if err != nil {
		return nil, empty_songs, err
	}

	if songs == nil {
		songs = empty_songs
	}

	return artist, songs, nil
}
