package database

import (
	"fmt"
	"log"
	"wp-backend/config"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() {
	err := mgm.SetDefaultConfig(nil, config.Config("DB"), options.Client().ApplyURI(fmt.Sprintf("mongodb://%v", config.Config("MONGO_URI"))))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connection opened to database")
}
