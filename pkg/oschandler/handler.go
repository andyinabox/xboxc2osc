package oschandler

import (
	"context"
	"path"

	"github.com/andyinabox/xboxcrelay/pkg/util"
	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
	"github.com/hypebeast/go-osc/osc"
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
		h.sendUint16s([]string{"ls"}, state.LeftStick[0], state.LeftStick[1])
		h.sendUint16s([]string{"ls", "x"}, state.LeftStick[0])
		h.sendUint16s([]string{"ls", "y"}, state.LeftStick[1])
	}
	if h.prev.RightStick != state.RightStick {
		h.sendUint16s([]string{"rs"}, state.RightStick[0], state.RightStick[1])
		h.sendUint16s([]string{"rs", "x"}, state.RightStick[0])
		h.sendUint16s([]string{"rs", "y"}, state.RightStick[1])
	}
	if h.prev.LeftTrigger != state.LeftTrigger {
		h.sendUint16s([]string{"lt"}, state.LeftTrigger)
	}
	if h.prev.RightTrigger != state.RightTrigger {
		h.sendUint16s([]string{"rt"}, state.RightTrigger)
	}
	if h.prev.DPad != state.DPad {
		if name, found := dpadStrMap[state.DPad]; found {
			h.sendBools([]string{"btn", name}, true)
		}
		if name, found := dpadStrMap[h.prev.DPad]; found {
			h.sendBools([]string{"btn", name}, false)
		}
	}
	if h.prev.MainButton != state.MainButton {
		if name, found := mainButtonsStrMap[state.MainButton]; found {
			h.sendBools([]string{"btn", name}, true)
		}
		if name, found := mainButtonsStrMap[h.prev.MainButton]; found {
			h.sendBools([]string{"btn", name}, false)
		}
	}
	if h.prev.SpecialButton != state.SpecialButton {
		if name, found := specialButtonsStrMap[state.SpecialButton]; found {
			h.sendBools([]string{"btn", name}, true)
		}
		if name, found := specialButtonsStrMap[h.prev.SpecialButton]; found {
			h.sendBools([]string{"btn", name}, false)
		}
	}

	return nil
}

func (h *Handler) sendUint16s(route []string, vals ...uint16) error {
	r := h.path(route)
	msg := osc.NewMessage(r)
	for _, v := range vals {
		msg.Append(util.NormalizeUint16(v))
	}
	return h.client.Send(msg)
}

func (h *Handler) sendBools(route []string, vals ...bool) error {
	r := h.path(route)
	msg := osc.NewMessage(r)
	for _, v := range vals {
		msg.Append(v)
	}
	return h.client.Send(msg)
}

func (h *Handler) path(s []string) string {
	args := []string{"/", h.conf.RoutePrefix}
	args = append(args, s...)
	return path.Join(args...)
}
