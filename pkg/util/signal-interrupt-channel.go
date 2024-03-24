package util

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalInterruptChannel() <-chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		close(interrupt)
	}()
	return interrupt
}
