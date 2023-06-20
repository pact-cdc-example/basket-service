package basket

import (
	"time"

	"github.com/pact-cdc-example/basket-service/pkg/cerr"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	requestTimeout = time.Second * 30
)

type Handler interface {
	SetupRoutes(fr fiber.Router)
}

type handler struct {
	service Service
	logger  *logrus.Logger
}

type NewHandlerOpts struct {
	S Service
	L *logrus.Logger
}

func NewHandler(opts *NewHandlerOpts) Handler {
	return &handler{
		service: opts.S,
		logger:  opts.L,
	}
}

func (h *handler) CreateBasket(c *fiber.Ctx) error {
	ctx := c.Context()

	var req CreateBasketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	basket, err := h.service.CreateBasket(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.JSON(basket)
}

func (h *handler) AddProductToBasket(c *fiber.Ctx) error {
	ctx := c.Context()

	basketID := c.Params("basket_id")

	var req AddProductToBasketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	req.BasketID = basketID

	basket, err := h.service.AddProductToBasket(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.JSON(basket)
}

func (h *handler) GetBasketByID(c *fiber.Ctx) error {
	ctx := c.Context()

	basketID := c.Params("basket_id")

	basket, err := h.service.GetBasketByID(ctx, basketID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.JSON(basket)
}

func (h *handler) AddBulkProductToBasket(c *fiber.Ctx) error {
	basketID := c.Params("basket_id")

	var req AddBulkProductToBasketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	req.BasketID = basketID

	basket, err := h.service.AddBulkProductToBasket(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.JSON(basket)

}

func (h *handler) SetupRoutes(fr fiber.Router) {
	basketGroup := fr.Group("/baskets")

	basketGroup.Post("/", h.CreateBasket)
	basketGroup.Post("/:basket_id", h.AddProductToBasket)
	basketGroup.Get("/:basket_id", h.GetBasketByID)
	basketGroup.Post("/:basket_id/bulk", h.AddBulkProductToBasket)
}
