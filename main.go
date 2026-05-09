package main

import (
	"Transaction/api/route"
	"Transaction/db"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	db.Initialize()
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message:": "H",
		})
	})
	route.SetupRoutes(app, db.DB)

	log.Fatal(app.Listen(":8080"))
}
