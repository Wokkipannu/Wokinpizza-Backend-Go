package handler

import (
	"log"
	"wp-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Topping struct {
	Toppings string `json:"toppings" xml:"toppings" form:"toppings"`
}

func GetToppings(c *fiber.Ctx) error {
	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching daily toppings"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Daily toppings", "data": toppings})
}

func UpdateToppings(c *fiber.Ctx) error {
	t := new(Topping)

	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching daily toppings"})
	}

	toppings.Toppings = t.Toppings
	err2 := mgm.Coll(toppings).Update(toppings)
	if err2 != nil {
		log.Printf("Failed to update toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error updating daily toppings"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Daily toppings updated", "data": toppings})
}
