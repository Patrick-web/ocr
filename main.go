package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	app.Get("/preview", func(c *fiber.Ctx) error {
		// Get the URL from the request query string
		url := c.Query("url")
		if url == "" {
			return c.Status(400).SendString("No URL provided")
		}

		// Generate the link preview
		preview, err := generateLinkPreview(c, url)
		if err != nil {
			return c.Status(500).SendString("Error generating link preview: " + err.Error())
		}

		// Return the link preview as a JSON response
		return c.JSON(preview)
	})
	println("Running")
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

// getAttr is a helper function that returns the value of a specified attribute
// for a given HTML node, or an empty string if the attribute does not exist.
func getAttr(n *html.Node, attr string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val, true
		}
	}
	return "", false
}

type LinkPreview struct {
	Title       string
	Description string
	Image       string
}

func generateLinkPreview(c *fiber.Ctx, url string) (*LinkPreview, error) {
	// Send a request to the URL and retrieve the HTML
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the HTML into a byte slice
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the HTML to extract the information we want
	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, err
	}

	var title, description, image string

	// Use a recursive function to traverse the HTML tree and extract the information we want
	var findInfo func(*html.Node)
	findInfo = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" {
				title = n.FirstChild.Data
			} else if n.Data == "meta" {
				nameAttr, nameOk := getAttr(n, "name")
				propAttr, propOk := getAttr(n, "property")
				if nameOk && nameAttr == "description" {
					description, _ = getAttr(n, "content")
				} else if propOk && propAttr == "og:image" {
					image, _ = getAttr(n, "content")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findInfo(c)
		}
	}
	findInfo(doc)

	// Return the link preview object
	return &LinkPreview{
		Title:       title,
		Description: description,
		Image:       image,
	}, nil
}
