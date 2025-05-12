package main

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
)

//Handlers
func getBooks(c *fiber.Ctx) error{
	return c.JSON(books)
}

// Handler functions
// getBooks godoc
// @Summary Get all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} Book
// @Router /book [get]
func getBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))//strconv คือ การสลับให้ จาก int เป็น string ในภาษา go
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, book := range books {//ปกติต้องเป็น index แต่ใน Go Lang ถ้าไม่ใช่ index ใช้ _เพื่อบอกว่ามองข้ามสิ่งที่เรียกว่า index ออกไป
		if book.ID == id {
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	book.ID = len(books) + 1
	books = append(books, *book)

	return c.JSON(book)
}

func updateBook(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == id {
			book.Title = bookUpdate.Title
			book.Author = bookUpdate.Author
			books[i] = book
			return c.JSON(book)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func deleteBook(c *fiber.Ctx) error{
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for i, book := range books{
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)//[1,2,3,4,5] -> [1,2] + [4,5] = [1,2,4,5]
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}