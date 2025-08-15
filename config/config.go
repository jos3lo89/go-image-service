// Package config: maneja la carga de las variables de entorno
package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        string `envconfig:"PORT" required:"true"`
	UploadDir   string `envconfig:"UPLOAD_DIR" required:"true"`
	URLClient   string `envconfig:"URL_CLIENT" required:"true"`
	MaxFileSize string `envconfig:"MAX_FILE_SIZE" required:"true"`
}

var AppConfig Config

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	err := envconfig.Process("", &AppConfig)
	if err != nil {
		log.Fatalf("Error processing environment variables: %s", err)
	}
}
