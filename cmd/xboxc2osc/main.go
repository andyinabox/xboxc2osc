package main

import (
	"github.com/andyinabox/xboxc2osc"
)

func main() {

	client := xboxc2osc.New(&xboxc2osc.Config{
		DeviceId:       [2]uint16{0x045e, 0x0b13},
		OscHostDomain:  "127.0.0.1",
		OscHostPort:    8000,
		OscRoutePrefix: "xboxc",
	})

	err := client.Run()
	if err != nil {
		panic(err)
	}
}
