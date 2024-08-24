package repositories

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"chords_app/internal/models"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Artist{}, &models.Song{}, &models.SongArtist{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateArtist(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	artist := models.Artist{Name: "New Artist"}

	err = repo.CreateArtist(&artist)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotZero(t, artist.ID, "expected artist ID to be set")
	assert.Equal(t, "New Artist", artist.Name, "expected artist name to be 'New Artist'")
}

func TestGetArtists(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	artist1 := models.Artist{Name: "Artist 1"}
	artist2 := models.Artist{Name: "Artist 2"}

	db.Create(&artist1)
	db.Create(&artist2)

	// Test case: Retrieve all artists
	artists, err := repo.GetArtists()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Len(t, *artists, 2, "expected 2 artists")
	assert.Equal(t, "Artist 1", (*artists)[0].Name, "expected first artist to be 'Artist 1'")
	assert.Equal(t, "Artist 2", (*artists)[1].Name, "expected second artist to be 'Artist 2'")
}

func TestGetArtistById(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	// Create test data
	artist := models.Artist{Name: "Artist By ID"}
	db.Create(&artist)

	// Test case: Retrieve artist by ID
	retrievedArtist, err := repo.GetArtistById(artist.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, artist.ID, retrievedArtist.ID, "expected artist ID to match")
	assert.Equal(t, "Artist By ID", retrievedArtist.Name, "expected artist name to be 'Artist By ID'")
}

func TestGetArtistById_NotFound(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	// Test case: Retrieve artist by non-existent ID
	retrievedArtist, err := repo.GetArtistById(999)
	assert.Nil(t, retrievedArtist, "expected no artist to be found")
	assert.NoError(t, err, "expected no error")
}

func TestGetArtistSongs_WithSongs(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	// Create test data
	artist := models.Artist{Name: "Test Artist"}
	song1 := models.Song{Title: "Song 1", Description: "Description 1", Content: "Content 1"}
	song2 := models.Song{Title: "Song 2", Description: "Description 2", Content: "Content 2"}

	db.Create(&artist)
	db.Create(&song1)
	db.Create(&song2)

	songArtist1 := models.SongArtist{ArtistID: artist.ID, SongID: song1.ID, TitleOrder: 1}
	songArtist2 := models.SongArtist{ArtistID: artist.ID, SongID: song2.ID, TitleOrder: 2}

	db.Create(&songArtist1)
	db.Create(&songArtist2)

	// Test case: Artist has songs
	songs, err := repo.GetArtistSongs(artist.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Len(t, songs, 2, "expected 2 songs")
	assert.Equal(t, "Song 1", songs[0].Title, "expected first song to be Song 1")
	assert.Equal(t, "Song 2", songs[1].Title, "expected second song to be Song 2")
}

func TestGetArtistSongs_WithoutSongs(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	// Create test data
	artistNoSongs := models.Artist{Name: "No Songs Artist"}
	db.Create(&artistNoSongs)

	// Test case: Artist has no songs
	songs, err := repo.GetArtistSongs(artistNoSongs.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Nil(t, songs, "expected no songs")
}

func TestGetArtistSongs_InvalidArtistID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	repo := NewGormArtistRepository(db)

	// Test case: Invalid artist ID
	songs, err := repo.GetArtistSongs(999)
	assert.Nil(t, songs, "expected no songs")
	assert.NoError(t, err, "expected no error")
}
