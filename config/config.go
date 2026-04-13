package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	KeyDB    KeyDBConfig
	JWT      JWTConfig
	CORS     CORSConfig
	MinIO    MinIOConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     string
	Host     string
	Timezone string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	MaxConns int
	MinConns int
}

type KeyDBConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	MaxRetries int
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	accessExpiryStr := viper.GetString("JWT_ACCESS_EXPIRY")
	accessExpiry, err := time.ParseDuration(accessExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY format: %w", err)
	}

	refreshExpiryStr := viper.GetString("JWT_REFRESH_EXPIRY")
	refreshExpiry, err := time.ParseDuration(refreshExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY format: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name:     viper.GetString("APP_NAME"),
			Env:      viper.GetString("APP_ENV"),
			Port:     viper.GetString("APP_PORT"),
			Host:     viper.GetString("APP_HOST"),
			Timezone: viper.GetString("APP_TIMEZONE"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Database: viper.GetString("POSTGRES_DATABASE"),
			MaxConns: viper.GetInt("POSTGRES_MAX_CONNS"),
			MinConns: viper.GetInt("POSTGRES_MIN_CONNS"),
		},
		KeyDB: KeyDBConfig{
			Host:       viper.GetString("KEYDB_HOST"),
			Port:       viper.GetString("KEYDB_PORT"),
			Password:   viper.GetString("KEYDB_PASSWORD"),
			DB:         viper.GetInt("KEYDB_DB"),
			MaxRetries: viper.GetInt("KEYDB_MAX_RETRIES"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
		},
		MinIO: MinIOConfig{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if cfg.Postgres.Host == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if cfg.KeyDB.Host == "" {
		return fmt.Errorf("KEYDB_HOST is required")
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}
