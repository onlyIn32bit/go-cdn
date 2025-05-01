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
	const maxFileSize = 20 * 1024 * 1024 // 20MB
	uploadDir := os.Getenv("UPLOAD_DIR")
	apiKey := os.Getenv("API_KEY")
	domain := os.Getenv("DOMAIN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Print(uploadDir)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatal("Cannot create uploads directory: ", err)
	}
	app := fiber.New()

	app.Post("/upload", func(c fiber.Ctx) error {
		if c.Get("X-API-Key") != apiKey {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized to upload file")
		}

		formFile, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "No image found")
		}

		if formFile.Size > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest, "File size limit exceeded: 20MB")
		}

		file, err := formFile.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot open image file")
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fiber.NewError(fiber.StatusUnsupportedMediaType, "Invalid file format")
		}

		// process image
		fileName := c.FormValue("name")
		if fileName == "" {
			fileName = uuid.New().String()
		}
		finalFileName := fmt.Sprintf("%s%s", fileName, ".jpeg")
		filePath := filepath.Join(uploadDir, finalFileName)

		resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)

		if err := imaging.Save(resizedImage, filePath, imaging.JPEGQuality(80)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save image")
		}

		info, err := os.Stat(filePath)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get file info")
		}

		return c.JSON(fiber.Map{
			"url":        fmt.Sprintf("%s/%s", domain, finalFileName),
			"size_after": info.Size(),
		})
	})
	app.Get("/*", static.New(uploadDir))
	app.Delete("/*", func(c fiber.Ctx) error {
		if c.Get("X-API-Key") != apiKey {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized to delete file")
		}

		fileName := c.Params("*")
		if fileName == "" {
			return fiber.NewError(fiber.StatusBadRequest, "No file name found")
		}

		filePath := filepath.Join(uploadDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "No file found")
		}

		if err := os.Remove(filePath); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot delete file")
		}

		return c.SendString("File deleted successfully")
	})

	app.Listen(fmt.Sprintf(":%s", port))
	log.Println("Server started on port", port)
}
