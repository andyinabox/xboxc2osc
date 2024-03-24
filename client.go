package xboxcrelay

import (
	"sync"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

type Client struct {
	mu       sync.Mutex
	client   *xboxc.Client
	current  *xboxc.State
	previous *xboxc.State
}

func NewClient() *Client {
	return &Client{
		current:  &xboxc.State{},
		previous: &xboxc.State{},
	}
}

func (c *Client) Open() error {
	return c.client.Open()
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Update() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// swap pointers
	ptr := c.previous
	c.previous = c.current
	c.current = ptr

	// update current state
	return c.client.Update(c.current)
}

func (c *Client) RightStick() (xy [2]float32, changed bool) {
	xy = [2]float32{
		uint16ToFloat32(c.current.RightStick[0]),
		uint16ToFloat32(c.current.RightStick[1]),
	}
	changed = c.current.RightStick != c.previous.RightStick
	return
}

func (c *Client) LeftStick() (xy [2]float32, changed bool) {
	xy = [2]float32{
		uint16ToFloat32(c.current.LeftStick[0]),
		uint16ToFloat32(c.current.LeftStick[1]),
	}
	changed = c.current.LeftStick != c.previous.LeftStick
	return
}

func (c *Client) RightTrigger() (v float32, changed bool) {
	v = uint16ToFloat32(c.current.RightTrigger)
	changed = c.current.RightTrigger != c.previous.RightTrigger
	return
}

func (c *Client) LeftTrigger() (v float32, changed bool) {
	v = uint16ToFloat32(c.current.LeftTrigger)
	changed = c.current.LeftTrigger != c.previous.LeftTrigger
	return
}

func (c *Client) DPad() (current DPadState, prev DPadState) {
	return DPadState(c.current.DPad), DPadState(c.previous.DPad)
}

func (c *Client) MainButton() (current MainButton, prev MainButton) {
	return MainButton(c.current.MainButton), MainButton(c.previous.MainButton)
}

func (c *Client) SpecialButton() (current SpecialButton, prev SpecialButton) {
	return SpecialButton(c.current.SpecialButton), SpecialButton(c.previous.SpecialButton)
}

func uint16ToFloat32(v uint16) float32 {
	return float32(v) / float32(65536-1) // -1 because of zero index
}
