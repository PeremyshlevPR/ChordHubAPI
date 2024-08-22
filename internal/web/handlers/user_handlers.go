package handlers

import (
	"chords_app/internal/config"
	"chords_app/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	service services.UserService
	roles   *config.Roles
}

func NewUserHandler(service services.UserService, roles *config.Roles) *UserHandler {
	return &UserHandler{service, roles}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password, h.roles.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.IssueAccessToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := h.service.IssueRefreshToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"userId":        fmt.Sprintf("%d", user.ID),
		"role":          user.Role,
		"email":         user.Email,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.IssueAccessToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := h.service.IssueRefreshToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"userId":        fmt.Sprintf("%d", user.ID),
		"role":          user.Role,
		"email":         user.Email,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newAccessToken, newRefreshToken, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
