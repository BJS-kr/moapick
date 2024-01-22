package main

import (
	"log"
	"moapick/db"
	"moapick/routes"

	_ "moapick/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @version 0.1.0
// @title moapick
// @description moapick 서비스의 api문서입니다.
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

	r.Get("/docs/*", swagger.HandlerDefault)
	routes.SetRouters(r)

	log.Fatal(r.Listen(":8080"))
}