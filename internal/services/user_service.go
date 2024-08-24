package services

import (
	"chords_app/internal/config"
	"chords_app/internal/models"
	r "chords_app/internal/repositories"
	"chords_app/internal/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(name, email, password, role string) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
	IssueAccessToken(userId uint, role, email string) (string, error)
	IssueRefreshToken(userId uint, role, email string) (string, error)
	Refresh(refreshToken string) (string, string, error)
}

type userService struct {
	repo      r.UserRepository
	jwtConfig *config.JWTConfig
}

func NewUserService(repo r.UserRepository, jwtConfig *config.JWTConfig) UserService {
	return &userService{repo, jwtConfig}
}

func (s *userService) Register(name, email, password, role string) (*models.User, error) {
	existingUser, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	auth_error := errors.New("invalid email or password")

	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return &models.User{}, auth_error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return &models.User{}, auth_error
	}

	return user, nil
}

func (s *userService) IssueAccessToken(userId uint, role, email string) (string, error) {
	accessTokenDuration := time.Duration(s.jwtConfig.AccessTokenExpTimeMin) * time.Minute
	return utils.IssueToken(userId, email, role, []byte(s.jwtConfig.AccessTokenSecretKey), accessTokenDuration)
}

func (s *userService) IssueRefreshToken(userId uint, role, email string) (string, error) {
	refreshTokenDuration := time.Duration(s.jwtConfig.RefreshTokenExpTimeDays) * (24 * time.Hour)
	return utils.IssueToken(userId, email, role, []byte(s.jwtConfig.RefreshTokenSecretKey), refreshTokenDuration)
}

func (s *userService) Refresh(refreshToken string) (string, string, error) {
	claims, err := utils.ValidateToken(refreshToken, []byte(s.jwtConfig.RefreshTokenSecretKey))
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.FindByEmail(claims.Email)
	if err != nil || user == nil {
		return "", "", errors.New("user not found")
	}

	accessToken, err := s.IssueAccessToken(user.ID, user.Role, user.Email)
	return accessToken, refreshToken, err
}
