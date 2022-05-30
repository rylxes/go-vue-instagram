package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"my-rest-api/handlers"
)

const port = ":8000"

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Get("/calculate/:account?", handlers.Calculate)

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
