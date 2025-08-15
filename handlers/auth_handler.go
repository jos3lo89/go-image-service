// Package handlers: contiene los handlers para la autenticación de usuarios.
package handlers

import (
	"jos3lo89/go-image-service/database"
	"jos3lo89/go-image-service/lib"
	"jos3lo89/go-image-service/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	input := new(LoginPayload)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := database.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
	}

	if !lib.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or password"})
	}

	token, err := lib.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	return c.Redirect("/")
}
