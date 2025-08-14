package main

import (
	"jos3lo89/go-image-service/config"
	"jos3lo89/go-image-service/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.Init()

	if err := os.MkdirAll(config.AppConfig.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Error al crear el directorio de subidas: %v", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/uploads", config.AppConfig.UploadDir)

	routes.SetupImageRoutes(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
