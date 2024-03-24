package main

import (
	"context"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/loghandler"
	"github.com/andyinabox/xboxcrelay/pkg/util"
)

func main() {
	interrupt := util.SignalInterruptChannel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	relay := xboxcrelay.New(loghandler.New(&loghandler.Config{
		MinWidth: 8,
		TabWidth: 0,
		Padding:  1,
		Padchar:  ' ',
		Flags:    0,
	}))

	go func() {
		err := relay.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	<-interrupt // wait for interrupt
}
