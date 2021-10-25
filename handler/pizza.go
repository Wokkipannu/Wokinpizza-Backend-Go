package handler

import (
	"log"
	"wp-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Pizza struct {
	ID          string `json:"id,omitempty" xml:"id,omitempty" form:"id,omitempty"`
	Name        string `json:"name" xml:"name" form:"name"`
	Description string `json:"description" xml:"description" form:"description"`
	Toppings    string `json:"toppings" xml:"toppings" form:"toppings"`
	Image       string `json:"image" xml:"image" form:"image"`
	Thumbnail   string `json:"thumbnail" xml:"thumbnail" form:"thumbnail"`
}

func GetAllPizza(c *fiber.Ctx) error {
	result := []models.Pizza{}

	err := mgm.Coll(&models.Pizza{}).SimpleFind(&result, bson.D{})
	if err != nil {
		log.Printf("Error finding all pizzas")
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching pizzas"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All pizzas", "data": result})
}

func GetPizza(c *fiber.Ctx) error {
	id := c.Params("id")

	pizza := &models.Pizza{}
	coll := mgm.Coll(pizza)

	err := coll.FindByID(id, pizza)
	if err != nil {
		log.Printf("Error finding pizza by id: %v", id)
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching pizza by id"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Pizza", "data": pizza})
}

func NewPizza(c *fiber.Ctx) error {
	p := new(Pizza)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	pizza := &models.Pizza{
		Name:        p.Name,
		Description: p.Description,
		Toppings:    p.Toppings,
		Image:       p.Image,
		Thumbnail:   p.Thumbnail,
	}

	err := mgm.Coll(pizza).Create(pizza)
	if err != nil {
		log.Printf("Error creating a pizza")
		return c.JSON(fiber.Map{"status": "error", "message": "Error creating a new pizza"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "New pizza", "data": pizza})
}

func DeletePizza(c *fiber.Ctx) error {
	id := c.Params("id")

	pizza := &models.Pizza{}
	coll := mgm.Coll(pizza)

	err := coll.FindByID(id, pizza)
	if err != nil {
		log.Printf("Error finding pizza by id: %v", id)
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching pizza by id"})
	}

	err2 := mgm.Coll(pizza).Delete(pizza)
	if err2 != nil {
		log.Printf("Error deleting pizza by id: %v", id)
		return c.JSON(fiber.Map{"status": "error", "message": "Error deleting pizza by id"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Pizza deleted"})
}

func UpdatePizza(c *fiber.Ctx) error {
	p := new(Pizza)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	pizza := &models.Pizza{}
	coll := mgm.Coll(pizza)

	err := coll.FindByID(p.ID, pizza)
	if err != nil {
		log.Printf("Error finding pizza by id: %v", p.ID)
		return c.JSON(fiber.Map{"status": "error", "message": "Error fetching pizza by id"})
	}

	pizza.Name = p.Name
	pizza.Description = p.Description
	pizza.Toppings = p.Toppings
	pizza.Image = p.Image
	pizza.Thumbnail = p.Thumbnail

	err2 := mgm.Coll(pizza).Update(pizza)
	if err2 != nil {
		log.Printf("Error updating pizza with id: %v", p.ID)
		return c.JSON(fiber.Map{"status": "error", "message": "Error updating pizza"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Pizza updated", "data": pizza})
}
