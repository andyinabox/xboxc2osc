package loghandler

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

// Config settings correspond to the `tabwriter.NewWriter` function
type Config struct {
	MinWidth int
	TabWidth int
	Padding  int
	Padchar  byte
	Flags    uint
}

type Handler struct {
	conf *Config
	w    *tabwriter.Writer
}

func New(conf *Config) *Handler {
	w := tabwriter.NewWriter(os.Stdout, conf.MinWidth, conf.TabWidth, conf.Padding, conf.Padchar, conf.Flags)
	return &Handler{conf, w}
}

func (h *Handler) HandleInput(ctx context.Context, state *xboxc.State) error {
	fmt.Fprintln(h.w, "LS X\tLS Y\tRS X\t RS Y\tLT\tRT\tDP\tMB\tSB\t")
	fmt.Fprintf(
		h.w,
		"%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		state.LeftStick[0],
		state.LeftStick[1],
		state.RightStick[0],
		state.RightStick[1],
		state.LeftTrigger,
		state.RightTrigger,
		state.DPad,
		state.MainButton,
		state.SpecialButton,
	)
	h.w.Flush()

	return nil
}
