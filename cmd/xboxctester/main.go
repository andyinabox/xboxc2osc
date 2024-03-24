package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

func connect(ctx context.Context, client *xboxc.Client, reqConnect <-chan struct{}) <-chan error {
	errors := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context expired, closing connect goroutine")
				close(errors)
				return
			case <-reqConnect:
				fmt.Println("Connecting...")
			connectLoop:
				for {
					fmt.Println("Attempting connection")
					err := client.Open()
					// if there's no error we can assume connection was successful and close out
					if err == nil {
						fmt.Println("Connection successful")
						errors <- nil
						break connectLoop
						// if there's an error, return error and let parent
						// decide whether to cancel
					} else {
						fmt.Println("Connection error: " + err.Error())
						errors <- err
					}
					// otherwise try again in 1 second
					fmt.Println("Device not found, retrying in 1 second")
					time.Sleep(time.Second)
				}
			}
		}
	}()
	return errors
}

func read(ctx context.Context, client *xboxc.Client, reqRead <-chan struct{}) <-chan error {
	errors := make(chan error)
	go func() {

		state := xboxc.NewState()
		w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.Debug)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context expired, closing read goroutine")
				close(errors)
				return
			case <-reqRead:
				fmt.Println("Starting read loop")
			readLoop:
				for {
					err := client.Update(state)
					if err != nil {
						errors <- err
						break readLoop
					}

					fmt.Fprintln(w, "LS X\tLS Y\tRS X\t RS Y\tLT\tRT\tDP\tMB\tSB\t")
					fmt.Fprintf(
						w,
						"%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
						state.LeftStick[0],
						state.LeftStick[1],
						state.RightStick[0],
						state.RightStick[1],
						state.LeftTrigger,
						state.RightTrigger,
						state.DPad,
						state.MainButton,
						state.SpecialButton,
					)
					w.Flush()
				}
			}
		}
	}()
	return errors
}

func main() {
	// create context with cancel func
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := xboxc.New()
	// state := xboxc.NewState()
	// w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.Debug)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	connectReqStream := make(chan struct{})
	connectionErrors := connect(ctx, client, connectReqStream)
	readReqStream := make(chan struct{})
	readErrors := read(ctx, client, readReqStream)

	connectReqStream <- struct{}{}

mainLoop:
	for {
		select {
		case <-interrupt:
			fmt.Println("Recieved interrupt, shutting down...")
			break mainLoop
		case err := <-connectionErrors:
			fmt.Println("Recieved error message from connection")
			if err == nil {
				fmt.Println("Connection was successful")
				readReqStream <- struct{}{}
			} else if err != xboxc.ErrDeviceNotFound {
				panic(err)
			}
		case err := <-readErrors:
			fmt.Println("Recieved error message from reader")
			if err != nil {
				if err == xboxc.ErrDisconnected {
					fmt.Println("Recieved disconnection error")
					connectReqStream <- struct{}{}
				} else {
					panic(err)
				}
			}
		}
	}
}
