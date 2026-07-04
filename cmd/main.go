package main

import (
	"ticket-system/config"
	"ticket-system/routes"
)

func main() {
	config.ConnectDatabase()

	router := routes.SetupRouter()

	router.Run(":8080")
}