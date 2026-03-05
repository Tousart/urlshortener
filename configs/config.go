package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	EnvPath    string    `yaml:"env_path"`
	Server     ServerCfg `yaml:"server"`
	PostgreSQL PostgreSQLCfg
}

type ServerCfg struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:""`
	Port int    `yaml:"port" env:"SERVER_PORT" env-default:"8000"`
}

type PostgreSQLCfg struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

type Flags struct {
	Repository string
	ConfigPath string
}

func ParseFlags() *Flags {
	repository := flag.String("repo", "inmemory", "repository type")
	cfgPath := flag.String("config", "", "path to config")
	flag.Parse()
	return &Flags{
		Repository: *repository,
		ConfigPath: *cfgPath,
	}
}

func LoadConfig(cfgPath string) (*Config, error) {
	var cfg Config
	var err error

	if cfgPath != "" {
		if _, err = os.Stat(cfgPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("config path not exists: %s", cfgPath)
		}
		if err = cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
			return nil, fmt.Errorf("config: LoadConfig error: %w", err)
		}
	}

	if cfg.EnvPath != "" {
		if _, err = os.Stat(cfg.EnvPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("env file not exists: %s", cfg.EnvPath)
		}
		if err = godotenv.Load(cfg.EnvPath); err != nil {
			return nil, fmt.Errorf("config: LoadConfig error: %w", err)
		}
	}

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
