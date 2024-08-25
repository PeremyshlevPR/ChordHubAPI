package handlers

import (
	"chords_app/internal/config"
	"chords_app/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SongHandler struct {
	service     services.SongService
	rolesConfig *config.Roles
	validate    *validator.Validate
}

func NewSongHandlers(service services.SongService, rolesConfig *config.Roles, validate *validator.Validate) *SongHandler {
	return &SongHandler{service, rolesConfig, validate}
}

func (h *SongHandler) UploadSong(c *gin.Context) {
	user, exists := GetUserModel(c)
	if !exists {
		return
	}

	var req struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
		Content     string `json:"content" validate:"required"`
		ArtistIds   []uint `json:"artistIds" validate:"required,min=1"`
	}
	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	song, _, err := h.service.UploadSong(req.Title, req.Description, req.Content, user.ID, req.ArtistIds)
	if err != nil {
		var statusCode int
		if err.Error() == "artist not found" {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"id":          song.ID,
			"title":       song.Title,
			"description": song.Description,
			"content":     song.Content,
			"artistIds":   req.ArtistIds,
		},
	)
}

func (h *SongHandler) GetSong(c *gin.Context) {
	songId, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song ID"})
		return
	}

	song, err := h.service.GetSongWithArtists(songId)
	if err != nil {
		var statusCode int
		if err.Error() == "song not found" {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	artistIds := make([]uint, 0, len(song.Artists))
	for _, artist := range song.Artists {
		artistIds = append(artistIds, artist.ID)
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"id":          song.ID,
			"title":       song.Title,
			"description": song.Description,
			"content":     song.Content,
			"uploadedBy":  song.UploadedBy,
			"artistIds":   artistIds,
		},
	)
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	songId, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song ID"})
		return
	}

	user, exists := GetUserModel(c)
	if !exists {
		return
	}

	var req struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description"`
		Content     string `json:"content" validate:"required"`
		ArtistIds   []uint `json:"artistIds" validate:"required,min=1"`
	}
	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	song, songArtists, err := h.service.UpdateSong(songId, req.Title, req.Description, req.Content, req.ArtistIds)
	if err != nil {
		var statusCode int
		if err.Error() == "song not found" || err.Error() == "artist not found" {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	if song.UploadedBy != user.ID && user.Role != h.rolesConfig.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin user or song owner can edit it"})
		return
	}

	artistIds := make([]uint, 0, len(*songArtists))
	for _, artist := range *songArtists {
		artistIds = append(artistIds, artist.ID)
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"id":          song.ID,
			"title":       song.Title,
			"description": song.Description,
			"content":     song.Content,
			"artistIds":   artistIds,
		},
	)
}
