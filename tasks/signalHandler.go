package tasks

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler is handling incoming OS signals.
type SignalHandler struct{}

// NewSignalHandler creates and returns a new SignalHandler.
func NewSignalHandler() *SignalHandler {
	return &SignalHandler{}
}

// Start runs the signal handler. It will run indefinitely until SIGINT or SIGTERM are received.
func (s *SignalHandler) Start(ctx context.Context) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signals:
	case <-ctx.Done():
	}
}
