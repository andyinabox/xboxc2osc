package main

import (
	"context"
	"flag"
	"log"

	"github.com/andyinabox/xboxcrelay"
	"github.com/andyinabox/xboxcrelay/pkg/midipub"
	"gitlab.com/gomidi/midi/v2"
)

func main() {

	port := flag.String("port", "VLC Rack", "Midi port name")
	channel := flag.Int("channel", 0, "MIDI channel to broadcast data on")
	list := flag.Bool("list", false, "List available ports")

	flag.Parse()

	publisher := midipub.New(&midipub.Config{
		MidiPort: *port,
		MidiChan: uint8(*channel),
	})

	if *list {
		ports := midi.GetOutPorts()
		log.Println("Listing available ports...")
		for i, p := range ports {
			log.Printf("- %d: %s \n", i, p.String())
		}
		return
	}

	relay := xboxcrelay.New(publisher)

	log.Println("Connecting...")
	relay.WaitUntilOpen(context.Background())

	log.Println("Starting relay...")
	for {
		err := relay.Update()
		if err != nil {
			// todo: attempt to reconnect when connection lost
			log.Printf("Error: %v", err)
			return
		}
	}
}
