package routes

import (
	"jos3lo89/go-image-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/upload", handlers.HandleUploadFile)
	api.Delete("/image/:filename", handlers.HandleDeleteFile)
	api.Get("/dowload", handlers.HandleDownloadAll)
	api.Get("/images", handlers.HandleListFiles)

}
