package main

import (
	"chords_app/internal/repositories"
	"chords_app/internal/services"
	"chords_app/internal/web"
	"chords_app/internal/web/handlers"
	"net/http"

	"chords_app/internal/config"
	"chords_app/internal/database"
	"log/slog"
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

	userRepo := repositories.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo, &cfg.JWTConfig)
	userHandler := handlers.NewUserHandler(userService, &cfg.Roles)

	router := web.SetupRouter(userHandler)
	slog.Info("Starting HTTP server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, router)
}
