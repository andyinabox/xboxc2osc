package xboxcrelay

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

type MainButton uint8

const (
	MainButtonA  MainButton = 1
	MainButtonB  MainButton = 2
	MainButtonX  MainButton = 8
	MainButtonY  MainButton = 16
	MainButtonLB MainButton = 64
	MainButtonRB            = 128
)

type SpecialButton uint8

const (
	SpecialButtonSelect     SpecialButton = 4
	SpecialButtonStart      SpecialButton = 8
	SpecialButtonLeftStick  SpecialButton = 32
	SpecialButtonRightStick SpecialButton = 64
)
