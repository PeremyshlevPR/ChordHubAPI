package middleware

import (
	"chords_app/internal/config"
	"chords_app/internal/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	jwtConfig *config.JWTConfig
}

func NewMiddleware(jwtconfig *config.JWTConfig) *authMiddleware {
	return &authMiddleware{jwtconfig}
}

func (a *authMiddleware) getCurrentUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			slog.Error("Authorization token was not provided")
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			slog.Error("Incorrect Format of Auth Token")
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Format of Authorization Token "})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(clientToken, []byte(a.jwtConfig.AccessTokenSecretKey))

		if err != nil {
			slog.Error("Invalid Token Signature")
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid of Expired Auth Token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Next()
	}
}
