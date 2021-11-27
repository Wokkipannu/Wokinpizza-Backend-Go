package handler

import (
	"log"
	"wp-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Dailytopping struct {
	Toppings string `json:"toppings" xml:"toppings" form:"toppings"`
}

func GetToppings(c *fiber.Ctx) error {
	toppings := []models.Topping{}

	err := mgm.Coll(&models.Topping{}).SimpleFind(&toppings, bson.D{})
	if err != nil {
		log.Printf("Failed to fetch toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching  toppings"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Toppings", "data": toppings})
}

func GetDailyToppings(c *fiber.Ctx) error {
	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch daily toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching daily toppings"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Daily toppings", "data": toppings})
}

func UpdateDailyToppings(c *fiber.Ctx) error {
	t := new(Dailytopping)

	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch daily toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching daily toppings"})
	}

	toppings.Toppings = t.Toppings
	err2 := mgm.Coll(toppings).Update(toppings)
	if err2 != nil {
		log.Printf("Failed to update daily toppings")
		return c.JSON(fiber.Map{"status": "error", "message": "Error updating daily toppings"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Daily toppings updated", "data": toppings})
}
