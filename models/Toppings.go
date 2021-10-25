package models

import (
	"github.com/kamva/mgm/v3"
)

type Dailytopping struct {
	mgm.DefaultModel `bson:",inline"`
	Toppings         string `json:"toppings" bson:"toppings"`
}
