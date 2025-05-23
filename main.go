package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/swagger"
	_ "github.com/leng/fiber-test/docs"
)

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book //Slice to store books

func checkMiddleware(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil{
		log.Fatal("load .env error")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	
	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	app.Post("/login", login)

	
	//JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	// CRUD routes
	app.Use(checkMiddleware)


	app.Get("/book", getBooks)
	app.Get("/book/:id", getBook)
	app.Post("/book", createBook)
	app.Put("/book/:id", updateBook)
	app.Delete("/book/:id", deleteBook)

	app.Post("/upload", uploadFile)
	app.Get("/test-html", testHTML)

	app.Get("/config", getEnv)

	app.Listen("localhost:3000")
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

type User  struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`

}
var memberUser = User{
	Email: "user@example.com",
	Password: "password123",
}

 func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	//ถ้า user and email ไม่ตรงกันจะไม่สามารถผ่านเข้ามาได้ ErrUnauthorized 
	if user.Email != memberUser.Email || user.Password != memberUser.Password{
		return fiber.ErrUnauthorized //ไม่อณุญาติจากฝั่งของ server
	}

	//Create token
	token := jwt.New(jwt.SigningMethodHS256)

	//Set claims
	claims := token.Claims.(jwt.MapClaims)//เก็บข้อมูลขอตัวที่เป็น token เอาไว้
	claims["email"]  =  user.Email
	claims["role"]   =  "admin"
	claims["exp"]    =  time.Now().Add(time.Hour * 72).Unix()
	
	//Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)

	}
	return c.JSON(fiber.Map{
		"message" : "Login success",
		"token"   : t,
	})
}