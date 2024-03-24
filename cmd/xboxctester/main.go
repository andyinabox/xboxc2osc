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

func connect(ctx context.Context, errStream chan<- error, clientStream <-chan *xboxc.Client) <-chan *xboxc.Client {
	connectedClientStream := make(chan *xboxc.Client)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(connectedClientStream)
				return
			case client := <-clientStream:
				fmt.Println("Connecting...")
				for {
					err := client.Open()

					if err == nil {
						fmt.Println("Connection successful")
						connectedClientStream <- client
						break
					} else {
						// any other errors should be reported
						if err != xboxc.ErrDeviceNotFound {
							errStream <- err
							break
						}

						// otherwise retry
						fmt.Println("Device not found, retrying in 1 second")
						time.Sleep(time.Second)
					}
				}
			}
		}
	}()
	return connectedClientStream
}

func read(ctx context.Context, errStream chan<- error, clientStream <-chan *xboxc.Client) <-chan *xboxc.State {
	stateStream := make(chan *xboxc.State)
	state := xboxc.NewState()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(stateStream)
				return
			case client := <-clientStream:
				for {
					err := client.Update(state)
					if err != nil {
						errStream <- err
						break
					}
					stateStream <- state
				}
			}
		}
	}()
	return stateStream
}

type inputHandler interface {
	HandleInput(ctx context.Context, state *xboxc.State) error
}

func delegate(ctx context.Context, errStream chan<- error, stateStream <-chan *xboxc.State, handlers ...inputHandler) {

	handlerStreams := make([]chan *xboxc.State, len(handlers))
	for i, handler := range handlers {
		// shadowing stateStream here specific to handler
		stateStream := make(chan *xboxc.State)
		handlerStreams[i] = stateStream

		// goroutine for individual handlers
		go func() {
			for {
				select {
				case <-ctx.Done():
					close(stateStream)
					return
				case state := <-stateStream:
					err := handler.HandleInput(ctx, state)
					if err != nil {
						errStream <- err
					}
				}
			}
		}()
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case state := <-stateStream:
				fmt.Println("Delegate streams")
				for _, stream := range handlerStreams {
					stream <- state
				}
			}
		}
	}()

}

type logger struct {
	w *tabwriter.Writer
}

func (l *logger) HandleInput(ctx context.Context, state *xboxc.State) error {
	fmt.Fprintln(l.w, "LS X\tLS Y\tRS X\t RS Y\tLT\tRT\tDP\tMB\tSB\t")
	fmt.Fprintf(
		l.w,
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
	l.w.Flush()

	return nil
}

func newLogger(w *tabwriter.Writer) *logger {
	// w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.Debug)
	return &logger{w}
}

func main() {
	// create context with cancel func
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := xboxc.New()

	logger1 := newLogger(tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.Debug))
	logger2 := newLogger(tabwriter.NewWriter(os.Stdout, 8, 0, 1, '.', tabwriter.Debug))
	logger3 := newLogger(tabwriter.NewWriter(os.Stdout, 8, 0, 1, '_', tabwriter.Debug))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	errs := make(chan error)
	clientStream := make(chan *xboxc.Client)

	stateStream := read(ctx, errs, connect(ctx, errs, clientStream))
	delegate(ctx, errs, stateStream, logger1, logger2, logger3)

	clientStream <- client

loop:
	for {
		select {
		case <-interrupt:
			fmt.Println("Recieved interrupt, shutting down...")
			break loop
		case err := <-errs:
			if err != nil {
				if err == xboxc.ErrDisconnected {
					fmt.Println("Recieved disconnection error")
					clientStream <- client
				} else {
					panic(err)
				}
			}
		}
	}
}
