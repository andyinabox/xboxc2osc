package main

import (
	"context"
	"flag"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/loghandler"
	"github.com/andyinabox/xboxcrelay/pkg/util"
)

func main() {

	var minWidth, tabWidth, padding int
	var padstr string

	flag.IntVar(&minWidth, "minwidth", 8, "Minimum column width")
	flag.IntVar(&tabWidth, "tabwidth", 0, "Tab separator width")
	flag.IntVar(&padding, "padding", 1, "Cell padding")
	flag.StringVar(&padstr, "padchar", " ", "Character to use for cell padding")

	flag.Parse()

	interrupt := util.SignalInterruptChannel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	relay := xboxcrelay.New(loghandler.New(&loghandler.Config{
		MinWidth: minWidth,
		TabWidth: tabWidth,
		Padding:  padding,
		Padchar:  []byte(padstr)[0],
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
