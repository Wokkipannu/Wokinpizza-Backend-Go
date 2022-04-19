package main

import (
	"log"
	"wp-backend/config"
	"wp-backend/database"
	"wp-backend/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: true,
	}))

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
