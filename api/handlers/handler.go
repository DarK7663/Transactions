package handlers

import (
	"Transaction/api/dto"
	"Transaction/db/repository"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	repo *repository.TransactionRepository
}

func NewUserHandler(repo *repository.TransactionRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data, err := h.repo.SearchUser(uint(id))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	data := new(dto.UpdateUser)
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	_, err = h.repo.UpdateUser(uint(id), data.Name, data.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var data dto.CreateUser

	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	task, err := h.repo.CreateUser(data.Name, data.Email, data.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": task,
	})
}

func (h *UserHandler) RegisterUser(c fiber.Ctx) error {
	var regUser dto.RegisterUser

	if err := c.Bind().Body(&regUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.repo.RegisterUser(regUser.Name, regUser.Email, regUser.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"User Created:": user,
	})
}

func (h *UserHandler) AuthenticateUser(c fiber.Ctx) error {
	var authUser dto.AuthenticateUser

	if err := c.Bind().Body(&authUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	_, err := h.repo.AuthenticateUser(authUser.Email, authUser.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Authenticated successfully",
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if _, err := h.repo.SearchUser(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	if err := h.repo.DeleteUser(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}
