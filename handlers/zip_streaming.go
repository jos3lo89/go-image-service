package handlers

import (
	"archive/zip"
	"bufio"
	"io"
	"jos3lo89/go-image-service/config"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func HandleDownloadAll(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", `attachment; filename="images-backup.zip"`)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer w.Flush()

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

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
			log.Printf("Error durante la transmisión del ZIP: %v", err)
		}
	})

	return nil
}
