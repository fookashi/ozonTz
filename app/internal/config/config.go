package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type databaseType string

const (
	inMemory databaseType = "inmemory"
	postgres databaseType = "postgres"
)

type DatabaseConfig interface {
	DSN() string
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

type InMemoryConfig struct{}

func (c InMemoryConfig) DSN() string {
	return "inmemory"
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type Config struct {
	Port   string         `env:"PORT"`
	DBType databaseType   `env:"DB_TYPE"`
	DB     DatabaseConfig `env:"-"`
	RedisConfig
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	switch cfg.DBType {
	case postgres:
		var pgConfig PostgresConfig
		if err := cleanenv.ReadEnv(&pgConfig); err != nil {
			return nil, fmt.Errorf("failed to load postgres config: %w", err)
		}
		cfg.DB = pgConfig
	case inMemory:
		cfg.DB = InMemoryConfig{}
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.DBType)
	}

	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return cfg
}
