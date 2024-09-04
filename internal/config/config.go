package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" validate:"required"`
	Server     Server     `yaml:"server" validate:"required"`
	DB         DB         `yaml:"db" validate:"required"`
	JWTConfig  JWTConfig  `yaml:"jwt" validate:"required"`
	Roles      Roles      `yaml:"roles" validate:"required"`
	Opensearch Opensearch `yaml:"opensearch" validate:"required"`
}

type Server struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type DB struct {
	Path string `yaml:"path" validate:"required"`
}

type JWTConfig struct {
	AccessTokenExpTimeMin uint   `yaml:"access_token_exp_time_min" validate:"required"`
	AccessTokenSecretKey  string `yaml:"access_token_secret_key" validate:"required"`

	RefreshTokenExpTimeDays uint   `yaml:"refresh_token_exp_time_days" validate:"required"`
	RefreshTokenSecretKey   string `yaml:"refresh_token_secret_key" validate:"required"`
}

type Roles struct {
	Admin string `yaml:"admin" validate:"required"`
	User  string `yaml:"user" validate:"required"`
}

type Opensearch struct {
	Addresses []string `yaml:"addresses" validate:"required"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	IndexName string   `yaml:"index_name"`
}

func SetupConfig() (*Config, error) {
	var config Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return &config, fmt.Errorf("required environment variable CONFIG_PATH is not set")
	}

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		return &config, err
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return &config, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}
