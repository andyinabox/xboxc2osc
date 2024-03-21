package xboxc2osc

import "encoding/binary"

type ControllerState struct {
	LeftStick     [2]uint16 // [x, y]
	RightStick    [2]uint16 // [x, y]
	LeftTrigger   uint16
	RightTrigger  uint16
	DPad          DPadState
	MainButton    MainButton
	SpecialButton SpecialButton
}

func (s *ControllerState) Assign(data []byte) {
	s.LeftStick[0] = binary.BigEndian.Uint16(data[1:3])
	s.LeftStick[1] = binary.BigEndian.Uint16(data[3:5])
	s.RightStick[0] = binary.BigEndian.Uint16(data[5:7])
	s.RightStick[1] = binary.BigEndian.Uint16(data[7:9])
	s.LeftTrigger = binary.BigEndian.Uint16(data[9:11])
	s.RightTrigger = binary.BigEndian.Uint16(data[11:13])
	s.DPad = DPadState(data[13])
	s.MainButton = MainButton(data[14])
	s.SpecialButton = SpecialButton(data[15])
}

func NewControllerState() *ControllerState {
	return &ControllerState{
		LeftStick:  [2]uint16{},
		RightStick: [2]uint16{},
	}
}
