package xboxc

import (
	"context"
	"errors"

	"github.com/sstallion/go-hid"
)

const (
	vid = 0x045e
	pid = 0x0b13
)

var ErrTimeout = errors.New("xboxcontroller timeout")
var ErrDisconnect = errors.New("xboxcontroller was disconnected")
var ErrNotOpen = errors.New("xboxcontroller device is not open")

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
func (c *Client) Open(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			d, err := hid.OpenFirst(vid, pid)
			if err == nil && d != nil {
				c.device = d
				return nil
			}
		}
	}
}

func (c *Client) Close() error {
	if c.device == nil {
		return ErrNotOpen
	}
	return c.device.Close()
}

func (c *Client) Update(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrTimeout
	default:
		return c.update()
	}
}

func (c *Client) update() error {
	if c.device == nil {
		return ErrNotOpen
	}

	_, err := c.device.Read(c.data)
	if err != nil {
		if err == hid.ErrTimeout {
			return ErrDisconnect
		}
		return err
	}

	c.diffState.Update(c.data)
	return nil
}

func (c *Client) State() *DiffState {
	return c.diffState
}
