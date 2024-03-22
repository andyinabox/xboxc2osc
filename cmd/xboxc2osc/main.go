package main

import (
	"context"
	"flag"
	"log"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/oscpub"
)

func main() {

	host := flag.String("host", "127.0.0.1", "Osc reciever host")
	port := flag.Int("port", 8000, "Osc reciever port")
	prefix := flag.String("prefix", "xboxc2osc", "Osc path prefix")

	flag.Parse()

	publisher := oscpub.New(&oscpub.Config{
		HostDomain:  *host,
		HostPort:    *port,
		RoutePrefix: *prefix,
	})

	relay := xboxcrelay.New(publisher)

	log.Println("Connecting...")
	relay.Open(context.Background())

	log.Println("Starting relay...")
	for {
		err := relay.Update(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
