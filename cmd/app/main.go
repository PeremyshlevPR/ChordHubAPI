package main

import (
	"chords_app/internal/adapters/opensearch"
	"chords_app/internal/config"
	"chords_app/internal/database"
	"chords_app/internal/repositories"
	"chords_app/internal/services"
	"chords_app/internal/web"
	"chords_app/internal/web/handlers"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

func main() {
	cfg, err := config.SetupConfig()
	if err != nil {
		slog.Error("error in reading config:", slog.String("error", err.Error()))
		return
	}
	slog.Info("Config initialized")

	db, err := database.SetupDatabase(&cfg.DB)
	if err != nil {
		slog.Error("error in setup database:", slog.String("error", err.Error()))
		return
	}
	database.AutoMigrate(db)
	slog.Info("Database initialized")

	validate := validator.New()

	userRepo := repositories.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo, &cfg.JWTConfig)
	userHandler := handlers.NewUserHandler(userService, &cfg.Roles)

	opensearchClient, err := opensearch.CreateOpenSearchClient(&cfg.Opensearch)
	if err != nil {
		slog.Error("Failed to initialize opensearch client", slog.String("error", err.Error()))
		return
	}
	opensrearchAdapter := opensearch.NewOpenSearchAdapter(opensearchClient, cfg.Opensearch.IndexName)
	artistRepo := repositories.NewGormArtistRepository(db)
	artistService := services.NewArtistService(artistRepo, opensrearchAdapter, db)
	artistHandler := handlers.NewArtistHandlers(artistService, validate)

	songRepo := repositories.NewGormSongRepository(db)
	songService := services.NewSongService(songRepo, artistRepo)
	songHandler := handlers.NewSongHandlers(songService, &cfg.Roles, validate)

	router := web.SetupRouter(userHandler, artistHandler, songHandler, userService, &cfg.Roles)

	slog.Info("Starting HTTP server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		slog.Error("error starting server:", slog.String("error", err.Error()))
	}
}
