package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
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
			return c.Status(403).SendString("Unauthorized to upload file")
		}

		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).SendString("No file found")
		}
		fileExt := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
		filePath := filepath.Join(uploadDir, newFileName)

		err = c.SaveFile(file, filePath)
		if err != nil {
			return c.Status(500).SendString("Server failed to save file")
		}

		url := fmt.Sprintf("%s/%s", domain, newFileName)
		return c.JSON(fiber.Map{
			"url": url,
		})
	})
	app.Get("/*", static.New(uploadDir))

	app.Listen(":8090")
}
