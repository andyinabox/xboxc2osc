package xboxcrelay

import (
	"context"
	"fmt"
	"time"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

type InputHandler interface {
	HandleInput(context.Context, *xboxc.State) error
}

type Relay struct {
	client   *xboxc.Client
	handlers []InputHandler
}

func New(handlers ...InputHandler) *Relay {
	return &Relay{xboxc.New(), handlers}
}

func (r *Relay) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error)
	defer close(errs)
	clientStream := make(chan *xboxc.Client)
	defer close(clientStream)
	stateStream := read(ctx, errs, connect(ctx, errs, clientStream))
	delegate(ctx, errs, stateStream, r.handlers...)

	clientStream <- r.client

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			if err != nil {
				if err == xboxc.ErrDisconnected {
					fmt.Println("Received disconnection error")
					clientStream <- r.client
				} else {
					return err
				}
			}
		}
	}

}

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

func delegate(ctx context.Context, errStream chan<- error, stateStream <-chan *xboxc.State, handlers ...InputHandler) {

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
				for _, stream := range handlerStreams {
					stream <- state
				}
			}
		}
	}()

}
