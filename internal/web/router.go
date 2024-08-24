package web

import (
	"chords_app/internal/web/handlers"

	"github.com/gin-gonic/gin"
)

const apiV1Prefix = "/api/v1"

func SetupRouter(userHandler *handlers.UserHandler, artistHandler *handlers.ArtistHandler) *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group(apiV1Prefix)

	apiRouter.POST("/register", userHandler.Register)
	apiRouter.POST("/login", userHandler.Login)
	apiRouter.POST("/refresh", userHandler.Refresh)

	artistRoutes := apiRouter.Group("/artists")
	{
		artistRoutes.POST("", artistHandler.CreateArtist)
		artistRoutes.GET("", artistHandler.GetArtists)
		artistRoutes.GET("/:id", artistHandler.GetArtistInformation)
		artistRoutes.PUT("/:id", artistHandler.UpdateArtist)
		artistRoutes.DELETE("/:id", artistHandler.DeleteArtist)
	}

	return r
}
