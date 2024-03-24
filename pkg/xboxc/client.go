package xboxc

import (
	"github.com/sstallion/go-hid"
)

const (
	vid = 0x045e
	pid = 0x0b13
)

type Client struct {
	device *hid.Device
	data   []byte
}

func New() *Client {
	return &Client{
		data: make([]byte, 17),
	}
}

// Open will continue to search for the device until it is found
func (c *Client) Open() error {
	d, err := hid.OpenFirst(vid, pid)
	if err != nil {
		if isErrDeviceNotFound(err) {
			return ErrDeviceNotFound
		}
		return err
	}
	c.device = d
	return nil
}

func (c *Client) Close() error {
	if c.device == nil {
		return ErrNotOpen
	}
	return c.device.Close()
}

func (c *Client) Update(state *State) error {
	if c.device == nil {
		return ErrNotOpen
	}

	_, err := c.device.Read(c.data)
	if err != nil {
		if isErrDisconnected(err) {
			return ErrDisconnected
		}
		return err
	}

	state.Assign(c.data)
	return nil
}
