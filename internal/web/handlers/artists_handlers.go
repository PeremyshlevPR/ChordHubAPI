package handlers

import (
	"chords_app/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ArtistHandler struct {
	service  services.ArtistService
	validate *validator.Validate
}

func NewArtistHandlers(service services.ArtistService, validate *validator.Validate) *ArtistHandler {
	return &ArtistHandler{service, validate}
}

func (h *ArtistHandler) GetArtists(c *gin.Context) {
	artists, err := h.service.GetArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if artists == nil || len(*artists) == 0 {
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	var response []map[string]interface{}
	for _, artist := range *artists {
		response = append(response, map[string]interface{}{
			"id":          artist.ID,
			"name":        artist.Name,
			"description": artist.Description,
			"imageUrl":    artist.ImageUrl,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ArtistHandler) GetArtistInformation(c *gin.Context) {
	artistId, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist ID"})
		return
	}

	artist, songs, err := h.service.GetArtistInformation(artistId)
	if err != nil {
		var code int
		if err.Error() == "artist not found" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          artist.ID,
		"name":        artist.Name,
		"description": artist.Description,
		"imageUrl":    artist.ImageUrl,
		"songs":       songs,
	})
}

func (h *ArtistHandler) CreateArtist(c *gin.Context) {
	// TODO: protect handler with auth middleware

	var req struct {
		Name        string `json:"name" validate:"required,min=2"`
		Description string `json:"description"`
		ImageUrl    string `json:"imageUrl"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	artist, err := h.service.CreateArtist(
		req.Name,
		req.Description,
		req.ImageUrl,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          artist.ID,
		"name":        artist.Name,
		"description": artist.Description,
		"imageUrl":    artist.ImageUrl,
	})
}

func (h *ArtistHandler) UpdateArtist(c *gin.Context) {
	// TODO: protect handler with auth middleware

	artistId, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageUrl    string `json:"imageUrl"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	artist, err := h.service.UpdateArtist(
		artistId,
		req.Name,
		req.Description,
		req.ImageUrl,
	)
	if err != nil {
		var code int
		if err.Error() == "artist not found" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          artist.ID,
		"name":        artist.Name,
		"description": artist.Description,
		"imageUrl":    artist.ImageUrl,
	})
}

func (h *ArtistHandler) DeleteArtist(c *gin.Context) {
	// TODO: protect handler with auth middleware

	artistId, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist ID"})
		return
	}

	err = h.service.DeleteArtist(artistId)
	if err != nil {
		var code int
		if err.Error() == "artist not found" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "artist deleted successfully"})
}
