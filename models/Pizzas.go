package models

import (
	"github.com/kamva/mgm/v3"
)

type Pizza struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Description      string `json:"description" bson:"description"`
	Toppings         string `json:"toppings" bson:"toppings"`
	Image            string `json:"image" bson:"image"`
	Thumbnail        string `json:"thumbnail" bson:"thumbnail"`
}
