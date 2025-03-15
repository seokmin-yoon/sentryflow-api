package main

import (
	"log"
	"sentryflow-api/app"
)

func main() {
	log.Println("Starting API Server...")
	app.Run()
}