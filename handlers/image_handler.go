package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"jos3lo89/go-image-service/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleUploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Imagen requerida."})
	}

	ext := filepath.Ext(file.Filename)

	uniqueFilename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), "upload", ext)

	dstPath := filepath.Join(config.AppConfig.UploadDir, uniqueFilename)

	if err := c.SaveFile(file, dstPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "No se pudo guardar el archivo."})
	}

	fileURL := fmt.Sprintf("%s/uploads/%s", c.BaseURL(), uniqueFilename)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Archivo subido con éxito", "url": fileURL, "id": uniqueFilename})
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
	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("Archivo '%s' eliminado con éxito.", filename)})
}

func HandleDownloadAll(c *fiber.Ctx) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(config.AppConfig.UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileToZip, err := os.Open(path)

		if err != nil {
			return err
		}

		defer fileToZip.Close()

		zipFile, err := zipWriter.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, fileToZip)
		return err
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error al crear el archivo ZIP."})
	}

	zipWriter.Close()

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", `attachment; filename="images-backup.zip"`)

	return c.Send(buf.Bytes())
}

func HandleListFiles(c *fiber.Ctx) error {
	files, err := os.ReadDir(config.AppConfig.UploadDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "No se pudo leer el directorio de imágenes."})
	}
	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return c.JSON(fiber.Map{"files": filenames})
}
