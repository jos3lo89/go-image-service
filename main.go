package main

import (
	"jos3lo89/go-image-service/config"
	"jos3lo89/go-image-service/database"
	"jos3lo89/go-image-service/models"
	"jos3lo89/go-image-service/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.Init()

	database.ConnectDB()

	// Crear usuario admin si no existe
	createAdminUser()

	if err := os.MkdirAll(config.AppConfig.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Error al crear el directorio de subidas: %v", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/public", "./public")
	app.Static("/uploads", config.AppConfig.UploadDir)

	routes.SetupImageRoutes(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}

// createAdminUser crea un usuario administrador por defecto al iniciar la app.
func createAdminUser() {
	var user models.User
	result := database.DB.Where("username = ?", "admin").First(&user)

	if result.RowsAffected == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		admin := models.User{Username: "admin", Password: string(hashedPassword)}
		database.DB.Create(&admin)
		log.Println("Usuario 'admin' creado con contraseña 'admin123'. ¡Cámbiala en producción!")
	}
}
