package oschandler

import (
	"context"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
	"github.com/hypebeast/go-osc/osc"
)

type Config struct {
	HostDomain  string
	HostPort    int
	RoutePrefix string
}

type Handler struct {
	conf   *Config
	client *osc.Client
	prev   xboxc.State
}

func New(conf *Config) *Handler {
	return &Handler{
		conf:   conf,
		client: osc.NewClient(conf.HostDomain, conf.HostPort),
	}
}

func (h *Handler) HandleInput(ctx context.Context, state *xboxc.State) error {
	// save a copy
	defer func() {
		h.prev = *state
	}()

	if h.prev.LeftStick != state.LeftStick {

	}
	if h.prev.RightStick != state.RightStick {

	}
	if h.prev.LeftTrigger != state.LeftTrigger {

	}
	if h.prev.RightTrigger != state.RightTrigger {

	}
	if h.prev.DPad != state.DPad {

	}
	if h.prev.MainButton != state.MainButton {

	}
	if h.prev.SpecialButton != state.SpecialButton {

	}
}
