package xboxcrelay

import (
	"context"
)

// var ErrDisconnected = errors.New("disconnected")

// EventPublisher is the destination for controller events.

type State interface {
	RightStick() (xy [2]float32, changed bool)
	LeftStick() (xy [2]float32, changed bool)
	RightTrigger() (v float32, changed bool)
	LeftTrigger() (v float32, changed bool)
	DPad() (current DPadState, previous DPadState)
	MainButton() (current MainButton, previous MainButton)
	SpecialButton() (current SpecialButton, pervious SpecialButton)
}

type Emitter interface {
	Emit(context.Context, State) error
}

type Relay struct {
	client   *Client
	emitters []Emitter
}

func New(emitters ...Emitter) *Relay {
	return &Relay{NewClient(), emitters}
}

func (r *Relay) Open() error {
	return r.client.Open()
}

func (r *Relay) Close() error {
	return r.client.Close()
}

func (r *Relay) Update(ctx context.Context) error {
	err := r.client.Update()
	for _, e := range r.emitters {
		err := e.Emit(ctx, r.client)
		if err != nil {
			return err
		}
	}
	return nil
}
