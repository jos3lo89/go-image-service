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

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AppConfig.URLClient,
	}))

	app.Use(logger.New())

	app.Static("/uploads", config.AppConfig.UploadDir)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "ok",
			"upload":   "/api/v1/upload",
			"download": "/api/v1/download",
			"list":     "/api/v1/images",
		})
	})
	routes.SetupImageRoutes(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
