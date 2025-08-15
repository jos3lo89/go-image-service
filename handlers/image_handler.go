// Package handlers: controladores de las imagenes
package handlers

import (
	"fmt"
	"jos3lo89/go-image-service/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ImageResponse struct {
	Image string `json:"image"`
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func HandleUploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Imagen requerida."})
	}

	// ext := filepath.Ext(file.Filename)

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Formato de archivo no permitido. Solo se aceptan: jpg, png, webp.",
		})
	}

	uniqueFilename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), "upload", ext)
	dstPath := filepath.Join(config.AppConfig.UploadDir, uniqueFilename)

	if err := c.SaveFile(file, dstPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "No se pudo guardar el archivo."})
	}

	fileURL := fmt.Sprintf("%s/uploads/%s", c.BaseURL(), uniqueFilename)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Archivo subido con éxito",
		"id":      uniqueFilename,
		"urlFull": fileURL,
		"url":     "/uploads/" + uniqueFilename,
	})
}

func HandleDeleteFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if strings.Contains(filename, "..") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Nombre de archivo no válido."})
	}
	filePath := filepath.Join(config.AppConfig.UploadDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "El archivo no existe."})
	}
	if err := os.Remove(filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Error al eliminar el archivo."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": fmt.Sprintf("Archivo '%s' eliminado con éxito.", filename)})
}

func HandleListFiles(c *fiber.Ctx) error {
	files, err := os.ReadDir(config.AppConfig.UploadDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "No se pudo leer el directorio de imágenes."})
	}

	var response []ImageResponse

	for _, file := range files {
		if !file.IsDir() {
			response = append(response, ImageResponse{Image: file.Name()})
		}
	}

	return c.JSON(response)
}
