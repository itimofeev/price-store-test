package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/itimofeev/price-store-test/internal/service"
)

func InitApp(srv *service.Service) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.DefaultErrorHandler(c, err)
		},
	})

	app.Use(logger.New())

	app.Post("/processCSV", processCSV(srv))
	app.Get("/listProducts", listProducts(srv))

	app.Get("/exampleCSV1", exampleCSV1())
	app.Get("/exampleCSV2", exampleCSV2())

	return app
}

func listProducts(srv *service.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		limitStr := c.Query("limit", "20")
		offsetStr := c.Query("offset", "0")
		orderStr := c.Query("order", "")

		limit := 20
		if parsed, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = int(parsed)
		}
		offset := 0
		if parsed, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			offset = int(parsed)
		}

		products, err := srv.ListProducts(c.Context(), orderStr, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to list products in service: %w", err)
		}
		return c.JSON(products)
	}
}

func processCSV(srv *service.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		url := c.Query("url", "")
		if url == "" {
			return errBadRequest("url is invalid")
		}

		if err := srv.ProcessCSV(c.Context(), url); err != nil {
			return fmt.Errorf("failed to process csv in service: %w", err)
		}
		return c.SendStatus(http.StatusOK)
	}
}

func exampleCSV1() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		example := `milk;10.02
honey;6.20
bread;7
`
		return c.SendString(example)
	}
}

func exampleCSV2() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		example := `milk;15.10
honey;6.20
butter;1.23
`
		return c.SendString(example)
	}
}

func errBadRequest(reason string) error {
	return fiber.NewError(http.StatusBadRequest, reason)
}
