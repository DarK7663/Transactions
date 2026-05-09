package route

import (
	"Transaction/api/handlers"
	"Transaction/db/repository"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewTaskRepository(db)

	handler := handlers.NewUserHandler(repo)

	api := app.Group("/api") // /api

	api.Get("/user/:id", handler.GetUser)

	api.Patch("/user/:id", handler.UpdateUser)

	api.Post("/user", handler.CreateUser)

	api.Post("/user/register", handler.RegisterUser)

	api.Post("user/authenticate", handler.AuthenticateUser)

	api.Post("user/delete/:id", handler.DeleteUser)
}
