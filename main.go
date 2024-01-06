package main

import (
	"fmt"
	"moapick/db"
	"moapick/routes"

	"github.com/otiai10/opengraph"
)

func main() {
	db.InitDB()
	r := routes.SetupRouter()
	og, err := opengraph.Fetch("https://naver.com")

	if err != nil {
		fmt.Print("error parsing og image")
	} else {
		if len(og.Image) > 0 {
			fmt.Print(og.Image[0].URL)
		}
	}

	r.Run(":8080")
}
