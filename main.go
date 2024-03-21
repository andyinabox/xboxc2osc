package xboxc2osc

import (
	"fmt"
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

	fmt.Println("Looking for device...")
	for {
		device, err = hid.OpenFirst(c.conf.DeviceId[0], c.conf.DeviceId[1])
		if err == nil && device != nil {
			break
		}
	}

	var msg *osc.Message
	data := make([]byte, 17)
	state := NewControllerState()

	for {

		_, err := device.Read(data)
		if err != nil {
			return err
		}

		state.Assign(data)
		fmt.Printf("%+v\n", state)

		// right stick
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

		msg = osc.NewMessage(c.oscPath("/ls/y"))
		msg.Append(uint16ToFloat32(state.RightStick[1] / 512))
		err = oscClient.Send(msg)
		if err != nil {
			return err
		}

		// left stick
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

		msg = osc.NewMessage(c.oscPath("/lt"))
		msg.Append(uint16ToFloat32(state.LeftTrigger / 512))
		err = oscClient.Send(msg)
		if err != nil {
			return err
		}

		msg = osc.NewMessage(c.oscPath("/rt"))
		msg.Append(uint16ToFloat32(state.RightTrigger / 512))
		err = oscClient.Send(msg)
		if err != nil {
			return err
		}

		dpad, found := DPadStateStrMap()[state.DPad]

		if found {

			var dpadx float32
			var dpady float32
			switch state.DPad {
			case DPadStateN:
				dpadx = 0.5
				dpady = 0.0
			case DPadStateNE:
				dpadx = 1.0
				dpady = 0.0
			case DPadStateE:
				dpadx = 1.0
				dpady = 0.5
			case DPadStateSE:
				dpadx = 1.0
				dpady = 1.0
			case DPadStateS:
				dpadx = 0.5
				dpady = 1.0
			case DPadStateSW:
				dpadx = 0.0
				dpady = 1.0
			case DPadStateW:
				dpadx = 0.0
				dpady = 0.5
			case DPadStateNW:
				dpadx = 0.0
				dpady = 0.0
			}

			msg = osc.NewMessage(c.oscPath("/dp/", dpad))
			msg.Append(true)
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

			msg = osc.NewMessage(c.oscPath("/dp"))
			msg.Append(dpadx)
			msg.Append(dpady)
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

		// special button
		specialButton, found := SpecialButtonStrMap()[state.SpecialButton]
		if found {
			msg = osc.NewMessage(c.oscPath("/btn/", specialButton))
			msg.Append(true)
			err = oscClient.Send(msg)
			if err != nil {
				return err
			}

			// we are also doplicating these under the ls/rs path
			// where other values for the joysticks are found
			if state.SpecialButton == SpecialButtonLeftStick {
				msg = osc.NewMessage(c.oscPath("/ls/btn"))
				msg.Append(true)
				err = oscClient.Send(msg)
				if err != nil {
					return err
				}
			}
			if state.SpecialButton == SpecialButtonRightStick {
				msg = osc.NewMessage(c.oscPath("/rs/btn"))
				msg.Append(true)
				err = oscClient.Send(msg)
				if err != nil {
					return err
				}
			}
		}

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
