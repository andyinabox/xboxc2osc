package xboxc2osc

import (
	"log"
	"path"

	"github.com/hypebeast/go-osc/osc"
	"github.com/sstallion/go-hid"
)

type Config struct {
	DeviceId       [2]uint16 // [VID, PID]
	OscHostDomain  string
	OscHostPort    int
	OscRoutePrefix string
}

type Client struct {
	conf *Config
}

func New(conf *Config) *Client {
	return &Client{conf}
}

func (c *Client) Run() error {
	var err error
	var device *hid.Device

	oscClient := osc.NewClient(c.conf.OscHostDomain, c.conf.OscHostPort)

	log.Println("Looking for device...")
	for {
		device, err = hid.OpenFirst(c.conf.DeviceId[0], c.conf.DeviceId[1])
		if err == nil && device != nil {
			break
		}
	}

	var msg *osc.Message
	data := make([]byte, 17)
	state := NewControllerState()
	prevState := NewControllerState()

	log.Println("Device connected, listening...")
	for {

		_, err := device.Read(data)
		if err != nil {
			return err
		}

		state.Assign(data)

		// right stick
		if state.RightStick != prevState.RightStick {
			msg = osc.NewMessage(c.oscPath("/rs"))
			msg.Append(uint16ToFloat32(state.RightStick[0]))
			msg.Append(uint16ToFloat32(state.RightStick[1]))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

			msg = osc.NewMessage(c.oscPath("/rs/x"))
			msg.Append(uint16ToFloat32(state.RightStick[0]))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

			msg = osc.NewMessage(c.oscPath("/rs/y"))
			msg.Append(uint16ToFloat32(state.RightStick[1] / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

		}

		// left stick
		if state.LeftStick != prevState.LeftStick {
			msg = osc.NewMessage(c.oscPath("/ls"))
			msg.Append(uint16ToFloat32(state.LeftStick[0] / 512))
			msg.Append(uint16ToFloat32(state.LeftStick[1] / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
			msg = osc.NewMessage(c.oscPath("/ls/x"))
			msg.Append(uint16ToFloat32(state.LeftStick[0] / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

			msg = osc.NewMessage(c.oscPath("/ls/y"))
			msg.Append(uint16ToFloat32(state.LeftStick[1] / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}

		// triggers

		if state.LeftTrigger != prevState.LeftTrigger {
			msg = osc.NewMessage(c.oscPath("/lt"))
			msg.Append(uint16ToFloat32(state.LeftTrigger / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}

		if state.RightTrigger != prevState.RightTrigger {
			msg = osc.NewMessage(c.oscPath("/rt"))
			msg.Append(uint16ToFloat32(state.RightTrigger / 512))
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}

		// dpad buttons
		for dpadStateType, str := range DPadStateStrMap() {
			msg = osc.NewMessage(c.oscPath("/dp/", str))
			msg.Append(dpadStateType == state.DPad)
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}

		// main button
		for mainButtonType, str := range MainButtonStrMap() {
			msg = osc.NewMessage(c.oscPath("/btn/", str))
			msg.Append(mainButtonType == state.MainButton)
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}
		// special buttons
		for specialButtonType, str := range SpecialButtonStrMap() {
			msg = osc.NewMessage(c.oscPath("/btn/", str))
			msg.Append(specialButtonType == state.SpecialButton)
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}
		}

		// swap pointers
		ptr := state
		state = prevState
		prevState = ptr
	}

}

func (c *Client) oscPath(p ...string) string {
	args := []string{"/", c.conf.OscRoutePrefix}
	args = append(args, p...)
	return path.Join(args...)
}

func uint16ToFloat32(v uint16) float32 {
	return float32(65536) / float32(v)
}
