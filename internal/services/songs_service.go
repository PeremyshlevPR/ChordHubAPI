package services

import (
	"chords_app/internal/models"
	"chords_app/internal/repositories"
	"errors"
)

type SongService interface {
	GetMostPopularSongs(period string, limit, offset uint) (*[]SongDTOWithViews, error)
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

type SongDTOWithViews struct {
	SongDTO
	Views uint
}

type songService struct {
	repo       repositories.SongRepository
	artistRepo repositories.ArtistRepository
}

func NewSongService(repo repositories.SongRepository, artistRepo repositories.ArtistRepository) SongService {
	return &songService{repo, artistRepo}
}

func (s *songService) GetMostPopularSongs(period string, limit, offset uint) (*[]SongDTOWithViews, error) {
	var days uint

	switch period {
	case "day":
		days = 1
	case "week":
		days = 7
	case "month":
		days = 30
	case "year":
		days = 365
	case "allTime":
		days = 0
	default:
		return nil, errors.New("invalid period, should by one of [day, week, month, year, allTime]")
	}

	songs, err := s.repo.GetPopularSongsForPeriod(days, limit, offset)
	if err != nil {
		return nil, err
	}
	return s.songstoSongDTO(songs)
}

func (s *songService) songstoSongDTO(songs *[]repositories.SongWithViews) (*[]SongDTOWithViews, error) {
	songDTOs := make([]SongDTOWithViews, 0, len(*songs))
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
		songDTOWithViews := SongDTOWithViews{
			songDTO,
			song.ViewCount,
		}
		songDTOs = append(songDTOs, songDTOWithViews)
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
	if err := s.repo.AddSongRequest(songId); err != nil {
		return nil, err
	}
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
