package web

import (
	"chords_app/internal/config"
	"chords_app/internal/services"
	"chords_app/internal/web/handlers"
	"chords_app/internal/web/middleware"

	"github.com/gin-gonic/gin"
)

const apiV1Prefix = "/api/v1"

func SetupRouter(
	userHandler *handlers.UserHandler,
	artistHandler *handlers.ArtistHandler,
	songHandler *handlers.SongHandler,
	userService services.UserService,
	rolesConfig *config.Roles,
) *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group(apiV1Prefix)
	apiRouter.POST("/register", userHandler.Register)
	apiRouter.POST("/login", userHandler.Login)
	apiRouter.POST("/refresh", userHandler.Refresh)
	apiRouter.GET("/artists", artistHandler.GetArtists)
	apiRouter.GET("/artists/:id", artistHandler.GetArtistInformation)
	apiRouter.GET("/songs/popular", songHandler.GetMostPopularSongs)
	apiRouter.GET("/songs/:id", songHandler.GetSong)

	authRequieredRouter := apiRouter.Group("/", middleware.AuthMiddleware(userService))
	authRequieredRouter.GET("/users/me", userHandler.GetUserInfo)
	authRequieredRouter.POST("/songs", songHandler.UploadSong)
	authRequieredRouter.PUT("songs/:id", songHandler.UpdateSong)

	adminOnlyRouter := authRequieredRouter.Group("/", middleware.VerifyRoleMiddleware(userService, rolesConfig.Admin))
	adminOnlyRouter.POST("/artists", artistHandler.CreateArtist)
	adminOnlyRouter.PUT("/artists/:id", artistHandler.UpdateArtist)
	adminOnlyRouter.DELETE("/artists/:id", artistHandler.DeleteArtist)
	adminOnlyRouter.POST("/users/create", userHandler.CreateNewUser)

	return r
}
