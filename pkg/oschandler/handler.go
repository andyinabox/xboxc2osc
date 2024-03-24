package oschandler

import (
	"context"
	"path"

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

	return nil
}

func (h *Handler) uint16(route []string, vals ...uint16) error {
	r := h.path(route)
	msg := osc.NewMessage(r)
	for _, v := range vals {
		msg.Append(v)
	}
	return h.client.Send(msg)
}

func (h *Handler) bool(route []string, vals ...bool) error {
	r := h.path(route)
	msg := osc.NewMessage(r)
	for _, v := range vals {
		msg.Append(v)
	}
	return h.client.Send(msg)
}

func (h *Handler) path(s []string) string {
	args := []string{"/", p.conf.RoutePrefix}
	args = append(args, s...)
	return path.Join(args...)
}

// func uint16ToFloat(n uint16) float32 {
// 	return
// }

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
