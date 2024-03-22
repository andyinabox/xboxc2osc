package xboxcrelay

import (
	"context"
	"errors"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

var ErrTimeout = errors.New("timeout")

// EventPublisher is the destination for controller events.
type EventPublisher interface {
	LeftStick(float32, float32) error
	RightStick(float32, float32) error
	LeftTrigger(float32) error
	RightTrigger(float32) error
	DPad(xboxc.DPadState, bool) error
	MainButton(xboxc.MainButton, bool) error
	SpecialButton(xboxc.SpecialButton, bool) error
}

type Relay struct {
	publisher  EventPublisher
	controller *xboxc.Client
}

func New(publisher EventPublisher) *Relay {
	return &Relay{publisher, xboxc.New()}
}

func (r *Relay) Open() error {
	return r.controller.Open()
}

func (r *Relay) WaitUntilOpen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			err := r.Open()
			if err == nil {
				return nil
			}
		}
	}
}

func (r *Relay) Update() error {
	err := r.controller.Update()
	if err != nil {
		return err
	}

	state := r.controller.State()

	// left joystick
	ls, changed := state.LeftStick()
	if changed {
		err := r.publisher.LeftStick(ls[0], ls[1])
		if err != nil {
			return err
		}
	}

	// right joystick
	rs, changed := state.RightStick()
	if changed {
		err := r.publisher.RightStick(rs[0], rs[1])
		if err != nil {
			return err
		}
	}

	// left trigger
	lt, changed := state.LeftTrigger()
	if changed {
		err := r.publisher.LeftTrigger(lt)
		if err != nil {
			return err
		}
	}

	// right trigger
	rt, changed := state.RightTrigger()
	if changed {
		err := r.publisher.RightTrigger(rt)
		if err != nil {
			return err
		}
	}

	// dpad
	currentDPad, prevDPad := state.DPad()
	if currentDPad != prevDPad {
		r.publisher.DPad(currentDPad, true)
		r.publisher.DPad(prevDPad, false)
	}

	// main buttons
	currentMainBtn, prevMainBtn := state.MainButton()
	if currentMainBtn != prevMainBtn {
		r.publisher.MainButton(currentMainBtn, true)
		r.publisher.MainButton(currentMainBtn, false)
	}

	// special buttons
	curSpecialBtn, prevSpecialBtn := state.SpecialButton()
	if curSpecialBtn != prevSpecialBtn {
		r.publisher.SpecialButton(curSpecialBtn, true)
		r.publisher.SpecialButton(prevSpecialBtn, false)
	}

	return nil
}
