package router

import (
	"manage-product-service/internal/handler"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

func NewRouter(productHandler *handler.ProductHandler) *fiber.App {
	app := fiber.New()

	app.Get("/docs/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/openapi.yaml")
	})
	app.Get("/api-docs/*", swagger.New(swagger.Config{
		URL: "/docs/openapi.yaml",
	}))

	app.Post("/product", productHandler.CreateProduct)
	app.Patch("/product/:id", productHandler.UpdateProduct)

	return app
}
