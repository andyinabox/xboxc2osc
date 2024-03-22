package xboxcrelay

import (
	"context"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

type EventPublisher interface {
	Float([]string, ...float32) error
	Bool([]string, ...bool) error
}

type Relay struct {
	publisher  EventPublisher
	controller *xboxc.Client
}

func New(publisher EventPublisher) *Relay {
	return &Relay{publisher, xboxc.New()}
}

func (r *Relay) Open(ctx context.Context) error {
	return r.controller.Open(ctx)
}

func (r *Relay) Update(ctx context.Context) error {
	err := r.controller.Update(ctx)
	if err != nil {
		return err
	}

	state := r.controller.State()

	// left joystick
	ls, changed := state.LeftStick()
	if changed {
		r.publisher.Float([]string{"ls"}, ls[0], ls[1])
		r.publisher.Float([]string{"ls", "x"}, ls[0])
		r.publisher.Float([]string{"ls", "y"}, ls[0])
	}

	// right joystick
	rs, changed := state.RightStick()
	if changed {
		r.publisher.Float([]string{"rs"}, rs[0], rs[1])
		r.publisher.Float([]string{"rs", "x"}, rs[0])
		r.publisher.Float([]string{"rs", "y"}, rs[0])
	}

	// left trigger
	lt, changed := state.LeftTrigger()
	if changed {
		r.publisher.Float([]string{"lt"}, lt)
	}

	// right trigger
	rt, changed := state.RightTrigger()
	if changed {
		r.publisher.Float([]string{"rt"}, rt)
	}

	// buttons
	var btnName string
	var found bool

	// dpad
	currentDPad, prevDPad := state.DPad()
	if currentDPad != prevDPad {
		if btnName, found = dpadStrMap[currentDPad]; found {
			r.publisher.Bool([]string{"dp", btnName}, true)
		}
		if btnName, found = dpadStrMap[prevDPad]; found {
			r.publisher.Bool([]string{"dp", btnName}, false)
		}
	}

	// main buttons
	currentMainBtn, prevMainBtn := state.MainButton()
	if currentMainBtn != prevMainBtn {
		if btnName, found = mainButtonsStrMap[currentMainBtn]; found {
			r.publisher.Bool([]string{"btn", btnName}, true)
		}
		if btnName, found = mainButtonsStrMap[prevMainBtn]; found {
			r.publisher.Bool([]string{"btn", btnName}, false)
		}
	}

	// special buttons
	curSpecialBtn, prevSpecialBtn := state.SpecialButton()
	if curSpecialBtn != prevSpecialBtn {
		if btnName, found = specialButtonsStrMap[curSpecialBtn]; found {
			r.publisher.Bool([]string{"btn", btnName}, true)
		}
		if btnName, found = specialButtonsStrMap[prevSpecialBtn]; found {
			r.publisher.Bool([]string{"btn", btnName}, false)
		}
	}

	return nil
}

var dpadStrMap = map[xboxc.DPadState]string{
	xboxc.DPadStateN:  "n",
	xboxc.DPadStateNE: "ne",
	xboxc.DPadStateE:  "e",
	xboxc.DPadStateSE: "se",
	xboxc.DPadStateS:  "s",
	xboxc.DPadStateSW: "sw",
	xboxc.DPadStateW:  "w",
	xboxc.DPadStateNW: "nw",
}
var mainButtonsStrMap = map[xboxc.MainButton]string{
	xboxc.MainButtonA:  "a",
	xboxc.MainButtonB:  "b",
	xboxc.MainButtonX:  "x",
	xboxc.MainButtonY:  "y",
	xboxc.MainButtonLB: "lb",
	xboxc.MainButtonRB: "rb",
}
var specialButtonsStrMap = map[xboxc.SpecialButton]string{
	xboxc.SpecialButtonSelect:     "select",
	xboxc.SpecialButtonStart:      "start",
	xboxc.SpecialButtonLeftStick:  "ls",
	xboxc.SpecialButtonRightStick: "rs",
}

// package xboxc2osc

// import (
// 	"context"
// 	"log"
// 	"path"

// 	"github.com/andyinabox/xboxc2osc/xboxc"
// 	"github.com/hypebeast/go-osc/osc"
// 	"github.com/sstallion/go-hid"
// )

// type EventPublisher interface {
// 	Float(context.Context, string, ...float32) error
// 	Bool(context.Context, string, ...bool) error
// }

// type Config struct {
// 	OscHostDomain  string
// 	OscHostPort    int
// 	OscRoutePrefix string
// }

// type Client struct {
// 	conf  *Config
// 	xboxc *xboxc.Client
// }

// func New(conf *Config) *Client {
// 	return &Client{conf}
// }

// func (c *Client) Run() error {
// 	var err error
// 	var device *hid.Device

// 	oscClient := osc.NewClient(c.conf.OscHostDomain, c.conf.OscHostPort)

// 	log.Println("Looking for device...")
// 	for {
// 		device, err = hid.OpenFirst(c.conf.DeviceId[0], c.conf.DeviceId[1])
// 		if err == nil && device != nil {
// 			break
// 		}
// 	}

// 	var msg *osc.Message
// 	data := make([]byte, 17)
// 	state := &ControllerState{}
// 	prevState := &ControllerState{}

// 	log.Println("Device connected, listening...")
// 	for {

// 		_, err := device.Read(data)
// 		if err != nil {
// 			if err == hid.ErrTimeout {
// 				log.Println("Device disconnected, exiting...")
// 				return nil
// 			}
// 			return err
// 		}

// 		state.Assign(data)

// 		// right stick
// 		if state.RightStick != prevState.RightStick {
// 			msg = osc.NewMessage(c.oscPath("/rs"))
// 			msg.Append(uint16ToFloat32(state.RightStick[0]))
// 			msg.Append(uint16ToFloat32(state.RightStick[1]))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}

// 			msg = osc.NewMessage(c.oscPath("/rs/x"))
// 			msg.Append(uint16ToFloat32(state.RightStick[0]))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}

// 			msg = osc.NewMessage(c.oscPath("/rs/y"))
// 			msg.Append(uint16ToFloat32(state.RightStick[1] / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}

// 		}

// 		// left stick
// 		if state.LeftStick != prevState.LeftStick {
// 			msg = osc.NewMessage(c.oscPath("/ls"))
// 			msg.Append(uint16ToFloat32(state.LeftStick[0] / 512))
// 			msg.Append(uint16ToFloat32(state.LeftStick[1] / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 			msg = osc.NewMessage(c.oscPath("/ls/x"))
// 			msg.Append(uint16ToFloat32(state.LeftStick[0] / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}

// 			msg = osc.NewMessage(c.oscPath("/ls/y"))
// 			msg.Append(uint16ToFloat32(state.LeftStick[1] / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// triggers

// 		if state.LeftTrigger != prevState.LeftTrigger {
// 			msg = osc.NewMessage(c.oscPath("/lt"))
// 			msg.Append(uint16ToFloat32(state.LeftTrigger / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if state.RightTrigger != prevState.RightTrigger {
// 			msg = osc.NewMessage(c.oscPath("/rt"))
// 			msg.Append(uint16ToFloat32(state.RightTrigger / 512))
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// dpad buttons
// 		for dpadStateType, str := range DPadStateStrMap() {
// 			msg = osc.NewMessage(c.oscPath("/dp/", str))
// 			msg.Append(dpadStateType == state.DPad)
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// main button
// 		for mainButtonType, str := range MainButtonStrMap() {
// 			msg = osc.NewMessage(c.oscPath("/btn/", str))
// 			msg.Append(mainButtonType == state.MainButton)
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		// special buttons
// 		for specialButtonType, str := range SpecialButtonStrMap() {
// 			msg = osc.NewMessage(c.oscPath("/btn/", str))
// 			msg.Append(specialButtonType == state.SpecialButton)
// 			err = oscClient.Send(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// swap pointers
// 		ptr := prevState
// 		prevState = state
// 		state = ptr
// 	}

// }

// func (c *Client) oscPath(p ...string) string {
// 	args := []string{"/", c.conf.OscRoutePrefix}
// 	args = append(args, p...)
// 	return path.Join(args...)
// }

// func uint16ToFloat32(v uint16) float32 {
// 	return float32(65536) / float32(v)
// }
