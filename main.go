package main

import (
	"japanism/db"
	"japanism/routes"
)

func main() {
	db.InitDB()
	r := routes.SetupRouter()

	r.Run(":8080")
}
