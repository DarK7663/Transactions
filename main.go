package main

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type UpdateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserCreate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	Initialize()
	app := fiber.New()
	repo := NewTaskRepository(db)

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:5500"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
	}))

	app.Get("/user/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		data, err := repo.SearchUser(uint(id))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"data": data,
		})
	})

	app.Patch("/user/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		data := new(UpdateUser)
		if err := c.Bind().Body(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		_, err = repo.UpdateUser(uint(id), data.Name, data.Email)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusNoContent).Send(nil)
	})

	app.Post("/user", func(c fiber.Ctx) error {
		var data UserCreate

		if err := c.Bind().Body(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		task, err := repo.CreateUser(data.Name, data.Email, data.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"data": task,
		})
	})

	app.Post("/user/register", func(c fiber.Ctx) error {
		var regUser RegisterUser

		if err := c.Bind().Body(&regUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		user, err := repo.RegisterUser(regUser.Name, regUser.Email, regUser.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"User Created:": user,
		})
	})

	app.Post("user/authenticate", func(c fiber.Ctx) error {
		var authUser AuthenticateUser

		if err := c.Bind().Body(&authUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		_, err := repo.AuthenticateUser(authUser.Email, authUser.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Authenticated successfully",
		})
	})

	app.Post("user/sendmoney", func(c fiber.Ctx) error {
		var transfer TransferRequest

		return nil
	})

	app.Delete("/user/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := repo.DeleteUser(uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusNoContent).Send(nil)
	})

	app.Listen("127.0.0.1:8080")
}
