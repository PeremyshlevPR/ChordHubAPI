package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string `yaml:"env"`
	Server    Server `yaml:"server"`
	DB        `yaml:"db"`
	JWTConfig `yaml:"jwt"`
	Roles     `yaml:"roles"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Path string `yaml:"path"`
}

type JWTConfig struct {
	AccessTokenExpTimeMin uint   `yaml:"access_token_exp_time_min"`
	AccessTokenSecretKey  string `yaml:"access_token_secret_key"`

	RefreshTokenExpTimeDays uint   `yaml:"refresh_token_exp_time_days"`
	RefreshTokenSecretKey   string `yaml:"refresh_token_secret_key"`
}

type Roles struct {
	Admin string `yaml:"admin"`
	User  string `yaml:"user"`
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
