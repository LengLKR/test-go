package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book //Slice to store books


func main() {
	if err := godotenv.Load(); err != nil{
		log.Fatal("load .env error")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	
	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	// CRUD routes
	app.Get("/book", getBooks)
	app.Get("/book/:id", getBook)
	app.Post("/book", createBook)
	app.Put("/book/:id", updateBook)
	app.Delete("/book/:id", deleteBook)

	app.Post("/upload",uploadFile)
	app.Get("/test-html",testHTML)

	app.Get("/config", getEnv)

	app.Listen(":3000")
}

func uploadFile(c *fiber.Ctx) error{

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/" + file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	
	return c.SendString("File upload complete!")

}

func testHTML(c *fiber.Ctx) error{
	return c.Render("index", fiber.Map{
		"Title" : "Hello, World!",
		"Name"  : "Leng",
	})
}


func getEnv(c *fiber.Ctx) error{
	secret := os.Getenv("SECRET")

	if secret == "" {
		secret = "defaultsecret"
	}

	return c.JSON(fiber.Map{
		"SECRET": os.Getenv("SECRET"),
	})
}