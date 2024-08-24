package handlers

import (
	"chords_app/internal/config"
	"chords_app/internal/models"
	"chords_app/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  services.UserService
	roles    *config.Roles
	validate *validator.Validate
}

func NewUserHandler(service services.UserService, roles *config.Roles) *UserHandler {
	return &UserHandler{
		service:  service,
		roles:    roles,
		validate: validator.New(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" validate:"required,min=2"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password, h.roles.User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.service.IssueAccessToken(user.ID, user.Role, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.service.IssueRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId":       user.ID,
		"role":         user.Role,
		"email":        user.Email,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.service.IssueAccessToken(user.ID, user.Role, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.service.IssueRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId":       user.ID,
		"role":         user.Role,
		"email":        user.Email,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	newAccessToken, newRefreshToken, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

func (h *UserHandler) CreateNewUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" validate:"required,min=2"`
		Email    string `json:"email" validate:"required,email"`
		Role     string `json:"role" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if !ValidateRequest(c, &req, h.validate) {
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": user.ID,
		"role":   user.Role,
		"email":  user.Email,
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not provided"})
		return
	}

	userModel := user.(*models.User)
	c.JSON(http.StatusOK, gin.H{
		"userId": userModel.ID,
		"role":   userModel.Role,
		"email":  userModel.Email,
	})
}
