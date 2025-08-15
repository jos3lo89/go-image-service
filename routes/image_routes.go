// Package routes: rutas de las imagenes
package routes

import (
	"jos3lo89/go-image-service/handlers"
	"jos3lo89/go-image-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(app *fiber.App) {
	api := app.Group("/api/v1", middleware.OptionalAuth())
	api.Post("/upload", handlers.HandleUploadFile)
	api.Delete("/image/:filename", handlers.HandleDeleteFile)
	api.Get("/download", handlers.HandleDownloadAll)
	api.Get("/images", handlers.HandleListFiles)
}
