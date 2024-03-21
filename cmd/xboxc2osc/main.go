package main

import (
	"flag"
	"strconv"

	"github.com/andyinabox/xboxc2osc"
)

func main() {

	vidStr := flag.String("vid", "045e", "HID vendor id as hexidecimal string")
	pidStr := flag.String("pid", "0b13", "HID product id as hexidecimal string")
	host := flag.String("host", "127.0.0.1", "Osc reciever host")
	port := flag.Int("port", 8000, "Osc reciever port")
	prefix := flag.String("prefix", "xboxc", "Osc path prefix")

	flag.Parse()

	vid, err := strconv.ParseUint(*vidStr, 16, 16)
	if err != nil {
		panic(err)
	}

	pid, err := strconv.ParseUint(*pidStr, 16, 16)
	if err != nil {
		panic(err)
	}

	client := xboxc2osc.New(&xboxc2osc.Config{
		DeviceId:       [2]uint16{uint16(vid), uint16(pid)},
		OscHostDomain:  *host,
		OscHostPort:    *port,
		OscRoutePrefix: *prefix,
	})

	err = client.Run()
	if err != nil {
		panic(err)
	}
}
