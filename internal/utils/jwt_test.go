package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestIssueToken(t *testing.T) {
	secretKey := []byte("my_secret_key")
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	expDuration := 15 * time.Minute

	tokenString, err := IssueToken(userID, email, role, secretKey, expDuration)

	assert.NoError(t, err, "IssueToken should not return an error")
	assert.NotEmpty(t, tokenString, "Generated token should not be empty")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	assert.NoError(t, err, "Parsing the token should not return an error")
	assert.True(t, token.Valid, "Token should be valid")
	assert.Equal(t, userID, claims.UserID, "UserID should match")
	assert.Equal(t, email, claims.Email, "Email should match")
	assert.Equal(t, role, claims.Role, "Role should match")
	assert.WithinDuration(t, time.Now().Add(expDuration), time.Unix(claims.ExpiresAt, 0), 1*time.Minute, "Expiration time should be within the expected range")
}

func TestValidateToken(t *testing.T) {
	secretKey := []byte("my_secret_key")
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	expDuration := 15 * time.Minute

	// First, issue a token
	tokenString, err := IssueToken(userID, email, role, secretKey, expDuration)
	assert.NoError(t, err, "IssueToken should not return an error")

	// Now validate the token
	claims, err := ValidateToken(tokenString, secretKey)
	assert.NoError(t, err, "ValidateToken should not return an error")
	assert.Equal(t, userID, claims.UserID, "UserID should match")
	assert.Equal(t, email, claims.Email, "Email should match")
	assert.Equal(t, role, claims.Role, "Role should match")
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secretKey := []byte("my_secret_key")
	invalidToken := "invalid.token.string"

	// Validate the invalid token
	_, err := ValidateToken(invalidToken, secretKey)
	assert.Error(t, err, "ValidateToken should return an error for an invalid token")
}
