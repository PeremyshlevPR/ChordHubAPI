package main

import (
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
}
