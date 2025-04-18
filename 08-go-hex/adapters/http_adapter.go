package adapters

import (
	"go-hex/core"

	"github.com/gofiber/fiber/v2"
)

type httpOrderHandler struct {
	service core.OrderService
}

func NewHttpOrderHandler(service core.OrderService) *httpOrderHandler {
	return &httpOrderHandler{service: service}
}

func (h *httpOrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order core.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.service.CreateOrder(order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
