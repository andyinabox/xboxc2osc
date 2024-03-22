package oscpub

import (
	"fmt"
	"log"
	"path"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
	"github.com/hypebeast/go-osc/osc"
)

type Config struct {
	HostDomain  string
	HostPort    int
	RoutePrefix string
}

type Publisher struct {
	conf   *Config
	client *osc.Client
}

func New(conf *Config) *Publisher {
	return &Publisher{conf, osc.NewClient(conf.HostDomain, conf.HostPort)}
}

func (p Publisher) Open() error {
	return nil
}

func (p *Publisher) LeftStick(x float32, y float32) error {
	p.sendFloats([]string{"ls"}, x, y)
	p.sendFloats([]string{"ls", "x"}, x)
	p.sendFloats([]string{"ls", "y"}, y)
	return nil
}

func (p *Publisher) RightStick(x float32, y float32) error {
	p.sendFloats([]string{"rs"}, x, y)
	p.sendFloats([]string{"rs", "x"}, x)
	p.sendFloats([]string{"rs", "y"}, y)
	return nil
}

func (p *Publisher) LeftTrigger(v float32) error {
	p.sendFloats([]string{"lt"}, v)
	return nil
}

func (p *Publisher) RightTrigger(v float32) error {
	p.sendFloats([]string{"rt"}, v)
	return nil
}

func (p *Publisher) DPad(b xboxc.DPadState, v bool) error {
	name, found := dpadStrMap[b]
	if !found {
		return fmt.Errorf("dpad button mapping not found for %v", b)
	}
	p.sendBools([]string{"db", name}, v)
	return nil
}

func (p *Publisher) MainButton(b xboxc.MainButton, v bool) error {
	name, found := mainButtonsStrMap[b]
	if !found {
		return fmt.Errorf("main button mapping not found for %v", b)
	}
	p.sendBools([]string{"btn", name}, v)
	return nil
}

func (p *Publisher) SpecialButton(b xboxc.SpecialButton, v bool) error {
	name, found := specialButtonsStrMap[b]
	if !found {
		return fmt.Errorf("special button mapping not found for %v", b)
	}
	p.sendBools([]string{"btn", name}, v)
	return nil
}

func (p *Publisher) sendFloats(route []string, floats ...float32) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	for _, f := range floats {
		msg.Append(f)
	}
	log.Printf("%s: %v", r, floats)
	return p.client.Send(msg)
}

func (p *Publisher) sendBools(route []string, bools ...bool) error {
	r := p.path(route)
	msg := osc.NewMessage(r)
	for _, b := range bools {
		msg.Append(b)
	}
	log.Printf("%s: %v", r, bools)
	return p.client.Send(msg)
}

func (p *Publisher) path(s []string) string {
	args := []string{"/", p.conf.RoutePrefix}
	args = append(args, s...)
	return path.Join(args...)
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
