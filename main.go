package xboxcrelay

import (
	"context"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

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
