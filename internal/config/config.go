package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env"`
	Server Server `yaml:"server"`
	DB     `yaml:"db"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Path string `yaml:"path"`
}

func SetupConfig() (*Config, error) {
	var config Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return &config, fmt.Errorf("required environment variable CONFIG_PATH is not set")
	}

	err := cleanenv.ReadConfig(configPath, &config)
	return &config, err
}
