package bootstrap

import "elProfessor/tasks"

// SignalHandler bootstraps the signal handler.
func SignalHandler() *tasks.SignalHandler {
	return tasks.NewSignalHandler()
}