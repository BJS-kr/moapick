package main

import (
	"log"
	"moapick/db"
	"moapick/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	envError := godotenv.Load("test.env")

	if envError != nil {
		panic("cannot load env")
	}

	db.InitDB()

	config := fiber.Config{
		Prefork: false,
	}

	r := fiber.New(config)

	r.Use(helmet.New())
	r.Use(cors.New())
	r.Use(logger.New())
	r.Use(recover.New())

	routes.SetRouters(r)

	log.Fatal(r.Listen(":8080"))
}
