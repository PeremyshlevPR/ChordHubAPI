package main

import (
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

	artistRepo := repositories.NewGormArtistRepository(db)
	artistService := services.NewArtistService(artistRepo)
	artistHandler := handlers.NewArtistHandlers(artistService, validate)

	router := web.SetupRouter(userHandler, artistHandler)

	slog.Info("Starting HTTP server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		slog.Error("error starting server:", slog.String("error", err.Error()))
	}
}
