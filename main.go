package main

import (
	"moapick/routes"
)

func main() {
	// db.InitDB()
	r := routes.SetupRouter()

	r.Run(":8080")
}
