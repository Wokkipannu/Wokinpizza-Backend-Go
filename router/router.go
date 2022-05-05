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
	api.Post("/pizza", handler.AuthLevelAdmin, handler.NewPizza)
	api.Put("/pizza", handler.AuthLevelAdmin, handler.UpdatePizza)
	api.Delete("/pizza", handler.AuthLevelAdmin, handler.DeletePizza)

	// Daily topping routes
	api.Get("/dailytoppings", handler.GetDailyToppings)
	api.Put("/dailytoppings", handler.UpdateDailyToppings)

	// Topping routes
	api.Get("/toppings", handler.GetToppings)
	api.Post("/topping", handler.AuthLevelAdmin, handler.NewTopping)
	api.Put("/topping", handler.AuthLevelAdmin, handler.UpdateTopping)
	api.Delete("/topping", handler.AuthLevelAdmin, handler.DeleteTopping)
	api.Get("/random", handler.GetRandomToppings)

	// Auth
	api.Get("/auth/login", handler.Login)
	api.Get("/auth/logout", handler.Logout)
	api.Get("/auth/callback", handler.Callback)
	api.Get("/auth/user", handler.GetUser)

	// Test routes
	api.Get("/auth/admin", handler.AuthLevelAdmin, func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
}
