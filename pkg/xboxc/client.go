package xboxc

import (
	"errors"

	"github.com/sstallion/go-hid"
)

const (
	vid = 0x045e
	pid = 0x0b13
)

var ErrNotOpen = errors.New("xboxc device is not open")

type Client struct {
	device    *hid.Device
	data      []byte
	diffState *DiffState
}

func New() *Client {
	return &Client{
		data:      make([]byte, 17),
		diffState: newDiffState(),
	}
}

// Open will continue to search for the device until it is found
func (c *Client) Open() error {
	d, err := hid.OpenFirst(vid, pid)
	if err != nil {
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

func (c *Client) Update() error {
	if c.device == nil {
		return ErrNotOpen
	}

	_, err := c.device.Read(c.data)
	if err != nil {
		return err
	}

	c.diffState.Update(c.data)
	return nil
}

func (c *Client) State() *DiffState {
	return c.diffState
}
