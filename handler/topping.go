package handler

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wp-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Dailytopping struct {
	Toppings string `json:"toppings" xml:"toppings" form:"toppings"`
}

type Topping struct {
	ID      string `json:"id,omitempty" xml:"id,omitempty" form:"id,omitempty"`
	Topping string `json:"topping" xml:"topping" form:"topping"`
}

type ToppingResponse struct {
	Text  string
	Count int64
}

func GetToppings(c *fiber.Ctx) error {
	toppings := []models.Topping{}

	err := mgm.Coll(&models.Topping{}).SimpleFind(&toppings, bson.D{})
	if err != nil {
		log.Printf("Failed to fetch toppings")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching toppings"})
	}
	return c.JSON(fiber.Map{"message": "Toppings", "data": toppings})
}

func NewTopping(c *fiber.Ctx) error {
	t := new(Topping)

	if err := c.BodyParser(t); err != nil {
		return err
	}

	topping := &models.Topping{
		Topping: t.Topping,
	}

	err := mgm.Coll(topping).Create(topping)
	if err != nil {
		log.Printf("Failed to create topping")
		return c.Status(500).JSON(fiber.Map{"message": "Error creating topping"})
	}

	return c.JSON(fiber.Map{"message": "Topping created", "data": topping})
}

func DeleteTopping(c *fiber.Ctx) error {
	id := c.Params("id")

	topping := &models.Topping{}
	coll := mgm.Coll(topping)

	err := coll.FindByID(id, topping)
	if err != nil {
		log.Printf("Failed to fetch topping")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching topping"})
	}

	err2 := coll.Delete(topping)
	if err2 != nil {
		log.Printf("Failed to delete topping")
		return c.Status(500).JSON(fiber.Map{"message": "Error deleting topping"})
	}

	return c.JSON(fiber.Map{"message": "Topping deleted"})
}

func UpdateTopping(c *fiber.Ctx) error {
	t := new(Topping)

	if err := c.BodyParser(t); err != nil {
		return err
	}

	topping := &models.Topping{}
	coll := mgm.Coll(topping)

	err := coll.FindByID(t.ID, topping)
	if err != nil {
		log.Printf("Failed to fetch topping")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching topping"})
	}

	topping.Topping = t.Topping

	err2 := coll.Update(topping)
	if err2 != nil {
		log.Printf("Failed to update topping")
		return c.Status(500).JSON(fiber.Map{"message": "Error updating topping"})
	}

	return c.JSON(fiber.Map{"message": "Topping updated", "data": topping})
}

func GetRandomToppings(c *fiber.Ctx) error {
	var amount int64
	if c.Params("amount") != "" {
		a, parseErr := strconv.ParseInt(c.Params("amount"), 10, 64)
		if parseErr != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid amount"})
		}
		amount = a
	} else if c.Query("amount") != "" {
		a, parseErr := strconv.ParseInt(c.Query("amount"), 10, 64)
		if parseErr != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid amount"})
		}
		amount = a
	} else {
		amount = 4
	}

	rand.Seed(time.Now().UnixNano())

	toppings := []models.Topping{}

	err := mgm.Coll(&models.Topping{}).SimpleFind(&toppings, bson.D{})
	if err != nil {
		log.Printf("Failed to fetch toppings")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching toppings"})
	}

	randomToppings := make([]ToppingResponse, len(toppings))
	for amount > 0 {
		randIndex := rand.Intn(len(toppings))
		countLimit := math.Ceil(float64(amount) / 12)
		count := rand.Int63n(int64(countLimit))
		if count == 0 {
			count += 1
		}
		randomToppings[randIndex].Count += count
		randomToppings[randIndex].Text = toppings[randIndex].Topping
		amount -= count
	}

	var output []string
	for _, v := range randomToppings {
		if v.Count > 1 {
			output = append(output, fmt.Sprintf("%vx %v", v.Count, v.Text))
		} else if v.Count == 1 {
			output = append(output, v.Text)
		}
	}

	return c.JSON(fiber.Map{"message": "Random pizza fetched", "data": strings.Join(output[:], ", ")})
}

func GetDailyToppings(c *fiber.Ctx) error {
	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch daily toppings")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching daily toppings"})
	}

	return c.JSON(fiber.Map{"message": "Daily toppings", "data": toppings})
}

func UpdateDailyToppings(c *fiber.Ctx) error {
	t := new(Dailytopping)

	toppings := &models.Dailytopping{}
	coll := mgm.Coll(toppings)

	err := coll.First(bson.M{}, toppings)
	if err != nil {
		log.Printf("Failed to fetch daily toppings")
		return c.Status(500).JSON(fiber.Map{"message": "Error fetching daily toppings"})
	}

	toppings.Toppings = t.Toppings
	err2 := mgm.Coll(toppings).Update(toppings)
	if err2 != nil {
		log.Printf("Failed to update daily toppings")
		return c.Status(500).JSON(fiber.Map{"message": "Error updating daily toppings"})
	}

	return c.JSON(fiber.Map{"message": "Daily toppings updated", "data": toppings})
}
