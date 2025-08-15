// Package handlers: contiene los handlers para la autenticación de usuarios.
package handlers

import (
	"jos3lo89/go-image-service/database"
	"jos3lo89/go-image-service/lib"
	"jos3lo89/go-image-service/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// LoginPayload con validación
type LoginPayload struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterPayload para registro de usuarios
type RegisterPayload struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

var validate = validator.New()

func Login(c *fiber.Ctx) error {
	input := new(LoginPayload)

	// Parsear el body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	// Validar los datos
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Buscar el usuario
	var user models.User
	result := database.DB.Where("username = ?", input.Username).First(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid username or password",
		})
	}

	// Verificar la contraseña
	if !lib.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid username or password",
		})
	}

	// Generar el token
	token, err := lib.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not generate token",
		})
	}

	// Establecer la cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		Secure:   true, // Usar HTTPS en producción
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	})
}

func Register(c *fiber.Ctx) error {
	input := new(RegisterPayload)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Verificar si el usuario ya existe
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "User already exists",
		})
	}

	// Hash de la contraseña
	hashedPassword, err := lib.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not process password",
		})
	}

	// Crear el usuario
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully",
	})
}

func Logout(c *fiber.Ctx) error {
	// Limpiar la cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func GetMe(c *fiber.Ctx) error {
	// Obtener el ID del usuario del contexto
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := database.DB.Select("id", "username", "email", "created_at").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user": user,
		},
	})
}

// RefreshToken renueva el token JWT
func RefreshToken(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	username := c.Locals("username").(string)

	// Generar nuevo token
	token, err := lib.GenerateToken(userID, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not generate token",
		})
	}

	// Actualizar la cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Token refreshed successfully",
	})
}
