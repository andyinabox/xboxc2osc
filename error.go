package xboxcrelay

import (
	"errors"
	"strings"
)

var ErrDisconnected = errors.New("device disconnected")

func IsErrDisconnected(err error) bool {
	return strings.Index(err.Error(), "hid_read_timeout") == 0
}
