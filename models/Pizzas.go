package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err2 := mgm.SetDefaultConfig(nil, os.Getenv("DB"), options.Client().ApplyURI(fmt.Sprintf("mongodb://%v", os.Getenv("MONGO_URI"))))
	if err2 != nil {
		log.Fatal(err2)
	}
}

type Pizza struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Description      string `json:"description" bson:"description"`
	Toppings         string `json:"toppings" bson:"toppings"`
	Image            string `json:"image" bson:"image"`
	Thumbnail        string `json:"thumbnail" bson:"thumbnail"`
}

func GetAllPizza() ([]Pizza, error) {
	result := []Pizza{}

	err := mgm.Coll(&Pizza{}).SimpleFind(&result, bson.D{})
	if err != nil {
		log.Printf("Error finding all pizzas")
		return result, err
	}

	return result, nil
}

func GetPizza(id string) (*Pizza, error) {
	pizza := &Pizza{}
	coll := mgm.Coll(pizza)

	err := coll.FindByID(id, pizza)
	if err != nil {
		log.Printf("Error finding pizza by id: %v", id)
		return pizza, err
	}

	return pizza, nil
}

func NewPizza(name string, description string, toppings string, image string, thumbnail string) *Pizza {
	return &Pizza{
		Name:        name,
		Description: description,
		Toppings:    toppings,
		Image:       image,
		Thumbnail:   thumbnail,
	}
}

func DeletePizza(id string) (string, error) {
	pizza, err := GetPizza(id)
	if err != nil {
		log.Printf(err.Error())
		result := fmt.Sprintf(err.Error())
		return result, err
	}

	err2 := mgm.Coll(pizza).Delete(pizza)
	if err2 != nil {
		log.Printf("Error deleting pizza by id %v", id)
		result := fmt.Sprintf("Error deleting pizza by id %v", id)
		return result, err
	}

	result := fmt.Sprintf("Pizza with id %v deleted", id)
	return result, nil
}

func UpdatePizza(id string, name string, description string, toppings string, image string, thumbnail string) (*Pizza, error) {
	pizza, err := GetPizza(id)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	pizza.Name = name
	pizza.Description = description
	pizza.Toppings = toppings
	pizza.Image = image
	pizza.Thumbnail = thumbnail

	err2 := mgm.Coll(pizza).Update(pizza)
	if err2 != nil {
		log.Printf("Failed to update pizza with id %v", id)
		return nil, err2
	}

	return pizza, nil
}
