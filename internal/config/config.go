package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// tag yaml опред. какое имя будет у соответствующего пар-ра в yaml файле
// tag env - имя пар-ра, если будем считывать его из переменной окружения
// tag env-default - значение по умолчанию для переменной окружения
// tag env-required - указывает, что переменная окружения обязательна для запуска приложения

type Config struct {
	Env string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	Storage StorageConfig `yaml:"storage" env-required:"true"`
	Postgres PostgresConfig `yaml:"postgres"`
	HTTPServer HTTPServerConfig `yaml:"http_server"`
}

type StorageConfig struct {
	Type string `yaml:"type" env-required:"true"`
}

type PostgresConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int `yaml:"port" env-required:"true"`
	User string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName string `yaml:"dbname" env-required:"true"`
	SSLMode string `yaml:"sslmode" env-required:"true"`
}

type HTTPServerConfig struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Приставка Must исп-ся, когда вместо ошибки уместно вызвать панику
func MustLoad() *Config{
	// Существует ли конфиг
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

// Собираем строку подкл. к бд
func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(p.User),
		url.QueryEscape(p.Password),
		p.Host,
		p.Port,
		p.DBName,
		p.SSLMode,
	)
}