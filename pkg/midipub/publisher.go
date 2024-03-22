package midipub

import (
	"fmt"
	"log"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type cc uint8

const (
	CCLeftStickX   cc = 0
	CCLeftStickY   cc = 1
	CCRightStickX  cc = 2
	CCRightStickY  cc = 3
	CCLeftTrigger  cc = 4
	CCRightTrigger cc = 5
)

var dpadNoteMap = map[xboxc.DPadState]uint8{
	xboxc.DPadStateN:  midi.C(0),
	xboxc.DPadStateNE: midi.Db(0),
	xboxc.DPadStateE:  midi.D(0),
	xboxc.DPadStateSE: midi.Eb(0),
	xboxc.DPadStateS:  midi.F(0),
	xboxc.DPadStateSW: midi.Gb(0),
	xboxc.DPadStateW:  midi.G(0),
	xboxc.DPadStateNW: midi.Ab(0),
}
var mainButtonsNoteMap = map[xboxc.MainButton]uint8{
	xboxc.MainButtonA:  midi.C(1),
	xboxc.MainButtonB:  midi.Db(1),
	xboxc.MainButtonX:  midi.D(1),
	xboxc.MainButtonY:  midi.Eb(1),
	xboxc.MainButtonLB: midi.F(1),
	xboxc.MainButtonRB: midi.Gb(1),
}
var specialButtonsNoteMap = map[xboxc.SpecialButton]uint8{
	xboxc.SpecialButtonSelect:     midi.G(1),
	xboxc.SpecialButtonStart:      midi.Ab(1),
	xboxc.SpecialButtonLeftStick:  midi.A(1),
	xboxc.SpecialButtonRightStick: midi.Bb(1),
}

type Config struct {
	MidiPort string
	MidiChan uint8
}

type Publisher struct {
	conf    *Config
	outPort *drivers.Out
	sendFn  func(msg midi.Message) error
}

func New(conf *Config) *Publisher {

	return &Publisher{
		conf: conf,
	}
}

func (p *Publisher) Open() error {
	port, err := drivers.OutByName(p.conf.MidiPort)
	if err != nil {
		return err
	}
	sendFn, err := midi.SendTo(port)
	if err != nil {
		return err
	}
	p.outPort = &port
	p.sendFn = sendFn
	return nil
}

func (p *Publisher) Close() error {
	p.outPort = nil
	p.sendFn = nil
	return nil
}

func (p *Publisher) LeftStick(x float32, y float32) error {
	err := p.sendCc(CCLeftStickX, float32ToUint8(x))
	if err != nil {
		return err
	}
	return p.sendCc(CCLeftStickY, float32ToUint8(y))
}

func (p *Publisher) RightStick(x float32, y float32) error {
	err := p.sendCc(CCRightStickX, float32ToUint8(x))
	if err != nil {
		return err
	}
	return p.sendCc(CCRightStickY, float32ToUint8(y))
}

func (p *Publisher) LeftTrigger(v float32) error {
	return p.sendCc(CCLeftTrigger, float32ToUint8(v))
}

func (p *Publisher) RightTrigger(v float32) error {
	return p.sendCc(CCRightTrigger, float32ToUint8(v))
}

func (p *Publisher) DPad(b xboxc.DPadState, v bool) error {
	note, found := dpadNoteMap[b]
	if !found {
		return fmt.Errorf("dpad button mapping not found for %v", b)
	}
	return p.sendNote(note, v)
}

func (p *Publisher) MainButton(b xboxc.MainButton, v bool) error {
	note, found := mainButtonsNoteMap[b]
	if !found {
		return fmt.Errorf("mmain button mapping not found for %v", b)
	}
	return p.sendNote(note, v)
}

func (p *Publisher) SpecialButton(b xboxc.SpecialButton, v bool) error {
	note, found := specialButtonsNoteMap[b]
	if !found {
		return fmt.Errorf("mmain button mapping not found for %v", b)
	}
	return p.sendNote(note, v)
}

func float32ToUint8(v float32) uint8 {
	f := v * 127
	return uint8(f)
}

func (p *Publisher) sendCc(controller cc, value uint8) error {
	log.Printf("cc: %v %v", controller, value)
	return p.sendFn(midi.ControlChange(p.conf.MidiChan, uint8(controller), value))
}

func (p *Publisher) sendNote(key uint8, b bool) error {
	if b {
		log.Printf("note on: %v", key)
		return p.sendFn(midi.NoteOn(p.conf.MidiChan, key, midi.On))
	} else {
		log.Printf("note off: %v", key)
		return p.sendFn(midi.NoteOff(p.conf.MidiChan, key))
	}
}
