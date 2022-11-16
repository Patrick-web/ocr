package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OCR Serve running")
	})
	app.Get("/imageUrl", func(c *fiber.Ctx) error {
		downloadImage(c.Query("url"), "./image.jpg")
		text := extractTextFromImage("./image.jpg")
		return c.SendString(text)
	})
	app.Get("/retry", func(c *fiber.Ctx) error {
		text := extractTextFromImage("./image.jpg")
		return c.SendString(text)
	})
	app.Listen(":3000")
}

// downloading image from url and saving it to the path as a jpg
func downloadImage(url string, path string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
}

//Extracting text from image using tesseract
func extractTextFromImage(path string) string {
	fmt.Println("Extracting text from image")
	client := gosseract.NewClient()
	client.Languages = []string{"eng"}
	client.SetWhitelist("0123456789+-?")
	defer client.Close()
	client.SetImage(path)
	text, _ := client.Text()
	fmt.Println(text)
	return text
}
