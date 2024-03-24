package main

import (
	"context"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/loghandler"
	"github.com/andyinabox/xboxcrelay/pkg/oschandler"
	"github.com/andyinabox/xboxcrelay/pkg/util"
)

func main() {
	interrupt := util.SignalInterruptChannel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := loghandler.New(&loghandler.Config{
		MinWidth: 8,
		TabWidth: 0,
		Padding:  1,
		Padchar:  ' ',
		Flags:    0,
	})

	osc := oschandler.New(&oschandler.Config{
		HostDomain:  "127.0.0.1",
		HostPort:    8000,
		RoutePrefix: "xboxcrelay",
	})

	relay := xboxcrelay.New(logger, osc)

	go func() {
		err := relay.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	<-interrupt // wait for interrupt
}
