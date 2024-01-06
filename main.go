package main

import (
	"moapick/db"
	"moapick/routes"

	"github.com/joho/godotenv"
)

func main() {
	envError := godotenv.Load("test.env")

	if envError != nil {
		panic("cannot load env")
	}

	db.InitDB()
	r := routes.SetupRouter()

	r.Run(":8080")
}
