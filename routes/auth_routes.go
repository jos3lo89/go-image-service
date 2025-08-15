package routes

import (
	"jos3lo89/go-image-service/handlers"
	"jos3lo89/go-image-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)
	auth.Post("/logout", handlers.Logout)

	api := app.Group("/api", middleware.Protected())
	api.Get("/me", handlers.GetMe)
	api.Post("/refresh", handlers.RefreshToken)
}
