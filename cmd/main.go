package main

import (
	"elProfessor/cmd/bootstrap"
	"elProfessor/cmd/config"
	"elProfessor/tasks"
)

func main() {
	config.Load()

	api := bootstrap.Api()

	tasks.RunTasks(api)
}