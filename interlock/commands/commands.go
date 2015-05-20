package commands

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	errChan = make(chan (error))
)

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	for _ = range sigChan {
		return
	}
}
