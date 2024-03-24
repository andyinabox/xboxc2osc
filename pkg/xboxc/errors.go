package xboxc

import (
	"errors"
	"strings"
)

var ErrNotOpen = errors.New("xboxc: client is not open")
var ErrDisconnected = errors.New("xboxc: device was disconnected")
var ErrDeviceNotFound = errors.New("xboxc: device not found")

func isErrDisconnected(err error) bool {
	return strings.Index(err.Error(), "hid_read_timeout") == 0
}

func isErrDeviceNotFound(err error) bool {
	return err.Error() == "No HID devices with requested VID/PID found in the system."
}
