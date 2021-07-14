package main

import (
	"elProfessor/cmd/bootstrap"
	"elProfessor/cmd/config"
	"elProfessor/internal/tasks"
	"log"
)

func main() {
	log.Println("Bootstrap initiated")

	config.Load()

	signalHandler := bootstrap.SignalHandler()
	db := bootstrap.Sqlite()
	api := bootstrap.Api(db)

	log.Println("Bootstrap finished. Heist API is starting")

	tasks.RunTasks(signalHandler, api)

	log.Println("Heist API finished gracefully")
}
