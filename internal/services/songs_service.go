package services

import (
	"chords_app/internal/models"
	"chords_app/internal/repositories"
	"errors"
)

type SongService interface {
	GetSongs(limit, offset uint, order_by string) (*[]SongDTO, error)
	UploadSong(title, description, content string, uploadedBy uint, artistIds []uint) (*models.Song, *[]models.SongArtist, error)
	UpdateSong(songId uint, title, description, content string, artistIds []uint) (*models.Song, *[]models.SongArtist, error)
	GetSongWithArtists(songId uint) (*models.Song, error)
	DeleteSong(songId uint) error
}

type SongDTO struct {
	ID      uint
	Title   string
	Artists []ArtistDTO
}

type songService struct {
	repo       repositories.SongRepository
	artistRepo repositories.ArtistRepository
}

func NewSongService(repo repositories.SongRepository, artistRepo repositories.ArtistRepository) SongService {
	return &songService{repo, artistRepo}
}

func (s *songService) GetSongs(limit, offset uint, order_by string) (*[]SongDTO, error) {
	songs, err := s.repo.GetSongs(limit, offset, order_by)
	if err != nil {
		return nil, err
	}

	songDTOs := make([]SongDTO, 0, len(*songs))
	for _, song := range *songs {
		artists := make([]ArtistDTO, 0, len(song.Artists))

		for _, songArtist := range song.Artists {
			artist, err := s.artistRepo.GetArtistById(songArtist.ArtistID)
			if err != nil || artist == nil {
				continue
			}
			artists = append(artists, ArtistDTO{artist.ID, artist.Name})
		}

		songDTO := SongDTO{
			ID:      song.ID,
			Title:   song.Title,
			Artists: artists,
		}
		songDTOs = append(songDTOs, songDTO)
	}

	return &songDTOs, nil
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
		artist, err := s.artistRepo.GetArtistById(artistId)
		if err != nil || artist == nil {
			return nil, nil, errors.New("artist not found")
		}

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

func (s *songService) UpdateSong(songId uint, title, description, content string, artistIds []uint) (*models.Song, *[]models.SongArtist, error) {
	song, err := s.repo.GetSongWithArtists(songId)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	if len(artistIds) > 0 {
		for _, songArtist := range song.Artists {
			s.repo.DeattachAuthor(&songArtist)
		}

		songArtists := make([]models.SongArtist, 0, len(artistIds))

		for i, artistId := range artistIds {
			artist, err := s.artistRepo.GetArtistById(artistId)
			if err != nil || artist == nil {
				return nil, nil, errors.New("artist not found")
			}

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
		song.Artists = songArtists
	}
	return song, &song.Artists, nil
}

func (s *songService) DeleteSong(songId uint) error {
	song, err := s.repo.GetSongById(songId)
	if err != nil || song == nil {
		return errors.New("song not found")
	}
	return s.repo.DeleteSong(song)
}
