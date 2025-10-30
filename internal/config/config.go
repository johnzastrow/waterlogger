package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Type    string         `yaml:"type"`
	SQLite  SQLiteConfig   `yaml:"sqlite"`
	MariaDB MariaDBConfig  `yaml:"mariadb"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type MariaDBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type AppConfig struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	SecretKey string `yaml:"secret_key"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Make path absolute
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save(configPath string) error {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(absPath, data, 0644)
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 2342,
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			SQLite: SQLiteConfig{
				Path: "waterlogger.db",
			},
			MariaDB: MariaDBConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "waterlogger",
				Password: "password",
				Database: "waterlogger",
			},
		},
		App: AppConfig{
			Name:      "Waterlogger",
			Version:   "1.0.0",
			SecretKey: "your-secret-key-change-this",
		},
	}
}