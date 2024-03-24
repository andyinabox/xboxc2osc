package main

import (
	"context"
	"flag"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/loghandler"
	"github.com/andyinabox/xboxcrelay/pkg/oschandler"
	"github.com/andyinabox/xboxcrelay/pkg/util"
)

func main() {
	var verbose bool
	var hostDomain, routePrefix string
	var hostPort int

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.StringVar(&hostDomain, "domain", "127.0.0.1", "OSC target host")
	flag.StringVar(&routePrefix, "prefix", "xboxcrelay", "OSC route prefix")
	flag.IntVar(&hostPort, "port", 8000, "OSC target port")
	flag.Parse()

	interrupt := util.SignalInterruptChannel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlers := []xboxcrelay.InputHandler{oschandler.New(&oschandler.Config{
		HostDomain:  "127.0.0.1",
		HostPort:    8000,
		RoutePrefix: "xboxcrelay",
	})}

	if verbose {
		logger := loghandler.New(&loghandler.Config{
			MinWidth: 8,
			TabWidth: 0,
			Padding:  1,
			Padchar:  ' ',
			Flags:    0,
		})
		handlers = append(handlers, logger)
	}

	relay := xboxcrelay.New(handlers...)

	go func() {
		err := relay.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	<-interrupt // wait for interrupt
}
