package xboxc2osc

type DPadState uint8

const (
	DPadStateN  DPadState = 1
	DPadStateNE DPadState = 2
	DPadStateE  DPadState = 3
	DPadStateSE DPadState = 4
	DPadStateS  DPadState = 5
	DPadStateSW DPadState = 6
	DPadStateW  DPadState = 7
	DPadStateNW DPadState = 8
)

func DPadStateStrMap() map[DPadState]string {
	return map[DPadState]string{
		DPadStateN:  "n",
		DPadStateNE: "ne",
		DPadStateE:  "e",
		DPadStateSE: "se",
		DPadStateS:  "s",
		DPadStateSW: "sw",
		DPadStateW:  "w",
		DPadStateNW: "nw",
	}
}

type MainButton uint8

const (
	MainButtonA  MainButton = 1
	MainButtonB  MainButton = 2
	MainButtonX  MainButton = 8
	MainButtonY  MainButton = 16
	MainButtonLB MainButton = 64
	MainButtonRB            = 128
)

func MainButtonStrMap() map[MainButton]string {
	return map[MainButton]string{
		MainButtonA:  "a",
		MainButtonB:  "b",
		MainButtonX:  "x",
		MainButtonY:  "y",
		MainButtonLB: "lb",
		MainButtonRB: "rb",
	}
}

type SpecialButton uint8

const (
	SpecialButtonSelect     SpecialButton = 4
	SpecialButtonStart      SpecialButton = 8
	SpecialButtonLeftStick  SpecialButton = 32
	SpecialButtonRightStick SpecialButton = 64
)

func SpecialButtonStrMap() map[SpecialButton]string {
	return map[SpecialButton]string{
		SpecialButtonSelect:     "select",
		SpecialButtonStart:      "start",
		SpecialButtonLeftStick:  "ls",
		SpecialButtonRightStick: "rs",
	}
}
