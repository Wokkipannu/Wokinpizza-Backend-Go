package router

import (
	"wp-backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Wokin Pizza API 🍕")
	})

	// Pizza routes
	api.Get("/pizza", handler.GetAllPizza)
	api.Get("/pizza/:id", handler.GetPizza)
	api.Post("/pizza", handler.NewPizza)
	api.Put("/pizza", handler.UpdatePizza)
	api.Delete("/pizza", handler.DeletePizza)

	// Daily topping routes
	api.Get("/dailytoppings", handler.GetDailyToppings)
	api.Put("/dailytoppings", handler.UpdateDailyToppings)

	// Topping routes
	api.Get("/toppings", handler.GetToppings)
	api.Get("/random", handler.GetRandomToppings)
}
