package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	MongoDB  MongoDBConfig
	Postgres PostgresConfig
	KeyDB    KeyDBConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
	Host string
}

type MongoDBConfig struct {
	URI         string
	Database    string
	MaxPoolSize int
	MinPoolSize int
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
	Host        string
	Port        string
	Password    string
	DB          int
	MaxRetries  int
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

var GlobalConfig *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	accessExpiry, err := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRY"))
	if err != nil {
		log.Fatalf("Invalid JWT_ACCESS_EXPIRY format: %v", err)
	}

	refreshExpiry, err := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))
	if err != nil {
		log.Fatalf("Invalid JWT_REFRESH_EXPIRY format: %v", err)
	}

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
			Host: viper.GetString("APP_HOST"),
		},
		MongoDB: MongoDBConfig{
			URI:         viper.GetString("MONGODB_URI"),
			Database:    viper.GetString("MONGODB_DATABASE"),
			MaxPoolSize: viper.GetInt("MONGODB_MAX_POOL_SIZE"),
			MinPoolSize: viper.GetInt("MONGODB_MIN_POOL_SIZE"),
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
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	GlobalConfig = config
	return config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if cfg.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
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