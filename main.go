package main

import (
	"moapick/db"
	"moapick/routes"
)

func main() {
	db.InitDB()
	r := routes.SetupRouter()

	r.Run(":8080")
}
