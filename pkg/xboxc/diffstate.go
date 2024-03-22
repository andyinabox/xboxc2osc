package xboxc

import (
	"sync"
)

type DiffState struct {
	mu       sync.Mutex
	current  *State
	previous *State
}

func newDiffState() *DiffState {
	return &DiffState{
		current:  &State{},
		previous: &State{},
	}
}

func (d *DiffState) Update(data []byte) {
	d.mu.Lock()
	// swap pointers
	ptr := d.previous
	d.previous = d.current
	d.current = ptr

	// update current from data
	d.current.Assign(data)
	d.mu.Unlock()
}

func (d *DiffState) RightStick() (xy [2]float32, changed bool) {
	xy = [2]float32{
		uint16ToFloat32(d.current.RightStick[0]),
		uint16ToFloat32(d.current.RightStick[1]),
	}
	changed = d.current.RightStick != d.previous.RightStick
	return
}

func (d *DiffState) LeftStick() (xy [2]float32, changed bool) {
	xy = [2]float32{
		uint16ToFloat32(d.current.LeftStick[0]),
		uint16ToFloat32(d.current.LeftStick[1]),
	}
	changed = d.current.LeftStick != d.previous.LeftStick
	return
}

func (d *DiffState) RightTrigger() (v float32, changed bool) {
	v = uint16ToFloat32(d.current.RightTrigger)
	changed = d.current.RightTrigger != d.previous.RightTrigger
	return
}

func (d *DiffState) LeftTrigger() (v float32, changed bool) {
	v = uint16ToFloat32(d.current.LeftTrigger)
	changed = d.current.LeftTrigger != d.previous.LeftTrigger
	return
}

func (d *DiffState) DPad() (current DPadState, prev DPadState) {
	return d.current.DPad, d.previous.DPad
}

func (d *DiffState) MainButton() (current MainButton, prev MainButton) {
	return d.current.MainButton, d.previous.MainButton
}

func (d *DiffState) SpecialButton() (current SpecialButton, prev SpecialButton) {
	return d.current.SpecialButton, d.previous.SpecialButton
}

func uint16ToFloat32(v uint16) float32 {
	return float32(v) / float32(65536-1) // -1 because of zero index
}
