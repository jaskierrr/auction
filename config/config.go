package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database   Database `envconfig:"db" required:"true"`
	ServerPort string   `envconfig:"serverport" required:"true" default:"8080"`
}

type Database struct {
	User     string `envconfig:"user" required:"true"`
	Password string `envconfig:"password" required:"true"`
	Host     string `envconfig:"host" required:"true"`
	Port     string `envconfig:"port" required:"true"`
	Name     string `envconfig:"name" required:"true"`
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal("Failed load envconfig " + err.Error())
	}

	return cfg
}
