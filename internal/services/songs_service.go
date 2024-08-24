package services

import (
	"chords_app/internal/models"
	"chords_app/internal/repositories"
	"errors"
)

type SongService interface {
}

type songService struct {
	repo repositories.SongRepository
}

func NewSongService(repo repositories.SongRepository) SongService {
	return &songService{repo}
}
func (s *songService) UploadSong(title, description, content string, uploadedBy uint, artistIds []uint) (*models.Song, *[]models.SongArtist, error) {
	song := models.Song{
		Title:       title,
		Description: description,
		Content:     content,
		UploadedBy:  uploadedBy,
	}

	if err := s.repo.CreateSong(&song); err != nil {
		return nil, nil, err
	}

	songArtists := make([]models.SongArtist, 0, len(artistIds))

	for i, artistId := range artistIds {
		songArtist := models.SongArtist{
			ArtistID:   artistId,
			SongID:     song.ID,
			TitleOrder: i,
		}
		if err := s.repo.AttachAuthor(&songArtist); err != nil {
			return nil, nil, err
		}
		songArtists = append(songArtists, songArtist)
	}

	return &song, &songArtists, nil
}

func (s *songService) GetSongWithArtists(songId uint) (*models.Song, error) {
	return s.repo.GetSongWithArtists(songId)
}

func (s *songService) UpdateSong(songId uint, title, description, content string) (*models.Song, error) {
	song, err := s.repo.GetSongById(songId)
	if err != nil {
		return nil, err
	}

	if song == nil {
		return nil, errors.New("song not found")
	}

	if title != "" {
		song.Title = title
	}
	if description != "" {
		song.Description = description
	}
	if content != "" {
		song.Content = content
	}

	if err := s.repo.UpdateSong(song); err != nil {
		return nil, err
	}
	return song, nil
}

func (s *songService) DeleteSong(songId uint) error {
	song, err := s.repo.GetSongById(songId)
	if err != nil {
		return err
	}
	if song == nil {
		return errors.New("song not found")
	}

	return s.repo.DeleteSong(song)
}
