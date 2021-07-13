package main

import (
	"elProfessor/cmd/bootstrap"
	"elProfessor/cmd/config"
	tasks "elProfessor/internal/tasks"
)

func main() {
	config.Load()

	api := bootstrap.Api()

	tasks.RunTasks(api)
}