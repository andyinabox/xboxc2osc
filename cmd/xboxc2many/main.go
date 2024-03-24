package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Running...")
	for {
		select {
		case <-interrupt:
			fmt.Println("Graceful shutdown...")
			os.Exit(0)
		default:
		}
	}

}
