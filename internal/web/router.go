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

	authRequieredRouter := apiRouter.Group("/", middleware.AuthMiddleware(userService))

	adminOnlyRouter := authRequieredRouter.Group("/", middleware.VerifyRoleMiddleware(userService, rolesConfig.Admin))
	adminOnlyRouter.POST("/artists", artistHandler.CreateArtist)
	adminOnlyRouter.PUT("/artists/:id", artistHandler.UpdateArtist)
	adminOnlyRouter.DELETE("/artists/:id", artistHandler.DeleteArtist)

	return r
}
