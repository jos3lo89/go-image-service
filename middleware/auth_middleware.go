// Package middleware: verifa la autenticidad de usuario
package middleware

import (
	"jos3lo89/go-image-service/config"
	"jos3lo89/go-image-service/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected es un middleware que verifica la validez del token JWT de la cookie.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el token de la cookie
		cookie := c.Cookies("jwt")

		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "No token provided",
			})
		}

		// Parsear y validar el token
		token, err := jwt.ParseWithClaims(cookie, &lib.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verificar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid or expired token",
			})
		}

		// Verificar que el token es válido
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid token",
			})
		}

		// Extraer los claims y guardarlos en el contexto
		if claims, ok := token.Claims.(*lib.JWTClaims); ok {
			// Guardar la información del usuario en el contexto
			c.Locals("userID", claims.UserID)
			c.Locals("username", claims.UserName)
		}

		return c.Next()
	}
}

func ProtectedWithRedirect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el token de la cookie
		cookie := c.Cookies("jwt")

		if cookie == "" {
			return c.Redirect("/login")
		}

		// Parsear y validar el token
		token, err := jwt.ParseWithClaims(cookie, &lib.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			// Si hay error o el token no es válido, redirigir al login
			return c.Redirect("/login")
		}

		// Extraer los claims y guardarlos en el contexto
		if claims, ok := token.Claims.(*lib.JWTClaims); ok {
			c.Locals("userID", claims.UserID)
			c.Locals("username", claims.UserName)
		}

		return c.Next()
	}
}

// GetUserFromContext obtiene la información del usuario del contexto
func GetUserFromContext(c *fiber.Ctx) (userID uint, username string, ok bool) {
	userIDInterface := c.Locals("userID")
	usernameInterface := c.Locals("username")

	if userIDInterface == nil || usernameInterface == nil {
		return 0, "", false
	}

	userID, ok1 := userIDInterface.(uint)
	username, ok2 := usernameInterface.(string)

	return userID, username, ok1 && ok2
}

// OptionalAuth verifica si hay un token válido pero no lo requiere
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")

		if cookie != "" {
			token, err := jwt.ParseWithClaims(cookie, &lib.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
				}
				return []byte(config.AppConfig.JWTSecret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*lib.JWTClaims); ok {
					c.Locals("userID", claims.UserID)
					c.Locals("username", claims.UserName)
					c.Locals("isAuthenticated", true)
				}
			}
		}

		// Si no hay token o es inválido, simplemente continúa sin autenticación
		if c.Locals("isAuthenticated") == nil {
			c.Locals("isAuthenticated", false)
		}

		return c.Next()
	}
}
