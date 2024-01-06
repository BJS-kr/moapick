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
		Prefork: true,
	}

	r := fiber.New(config)
	
  r.Use(helmet.New())
	r.Use(cors.New())
	r.Use(logger.New())
	r.Use(recover.New())
	routes.SetupRouter(r)
	// TODO 배포할 때 trusted proxy 설정
	// r.SetTrustedProxies([]string{})

	log.Fatal(r.Listen(":8080"))
}
