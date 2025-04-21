package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	const maxFileSize = 20 * 1024 * 1024
	uploadDir := os.Getenv("UPLOAD_DIR")
	apiKey := os.Getenv("API_KEY")
	domain := os.Getenv("DOMAIN")
	fmt.Print(uploadDir)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Post("/upload", func(c fiber.Ctx) error {
		if c.Get("X-API-Key") != apiKey {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized to upload file")
		}

		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file found")
		}

		if file.Size > maxFileSize {
			return c.Status(fiber.StatusBadRequest).SendString("File size limit exceeded")
		}

		fileSrc, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("No file found")
		}
		defer fileSrc.Close()

		img, _, err := image.Decode(fileSrc)
		if err != nil {
			return fiber.NewError(fiber.StatusUnsupportedMediaType, "Invalid image format")
		}

		dstImg := imaging.Resize(img, 1280, 0, imaging.Lanczos)

		fileExt := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
		filePath := filepath.Join(uploadDir, newFileName)

		err = imaging.Save(dstImg, filePath, imaging.JPEGQuality(80))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save image")
		}

		url := fmt.Sprintf("%s/%s", domain, newFileName)
		return c.JSON(fiber.Map{
			"url": url,
		})
	})
	app.Get("/*", static.New(uploadDir))
	app.Delete("/*", func(c fiber.Ctx) error {
		if c.Get("X-API-Key") != apiKey {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized to delete file")
		}

		fileName := c.Params("*")
		if fileName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("No file name found")
		}

		filePath := filepath.Join(uploadDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("No file found")
		}

		if err := os.Remove(filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot delete file")
		}

		return c.SendString("File deleted successfully")
	})

	app.Listen(":8090")
}
