package main

import (
	"log"
	"models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/kamva/mgm/v3"
)

type Pizza struct {
	ID          string `json:"id,omitempty" xml:"id,omitempty" form:"id,omitempty"`
	Name        string `json:"name" xml:"name" form:"name"`
	Description string `json:"description" xml:"description" form:"description"`
	Toppings    string `json:"toppings" xml:"toppings" form:"toppings"`
	Image       string `json:"image" xml:"image" form:"image"`
	Thumbnail   string `json:"thumbnail" xml:"thumbnail" form:"thumbnail"`
}

func main() {
	app := fiber.New()

	app.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	app.Use(cors.New())

	file, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip}:${port} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Europe/Helsinki",
		Output:     file,
	}))

	app.Get("/dashboard", monitor.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Wokin Pizza API 🍕")
	})

	app.Get("/pizza", func(c *fiber.Ctx) error {
		pizza, err := models.GetAllPizza()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.JSON(pizza)
	})

	app.Get("/pizza/:id", func(c *fiber.Ctx) error {
		pizza, err := models.GetPizza(c.Params("id"))
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.JSON(pizza)
	})

	app.Post("/pizza", func(c *fiber.Ctx) error {
		p := new(Pizza)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		pizza := models.NewPizza(p.Name, p.Description, p.Toppings, p.Image, p.Thumbnail)
		err := mgm.Coll(pizza).Create(pizza)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(pizza)
	})

	app.Delete("/pizza", func(c *fiber.Ctx) error {
		p := new(Pizza)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		res, err := models.DeletePizza(p.ID)
		if err != nil {
			return err
		}

		return c.SendString(res)
	})

	app.Put("/pizza", func(c *fiber.Ctx) error {
		p := new(Pizza)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		pizza, err := models.UpdatePizza(p.ID, p.Name, p.Description, p.Toppings, p.Image, p.Thumbnail)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(pizza)
	})

	app.Listen(":3000")
}
