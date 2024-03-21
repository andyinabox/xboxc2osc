# `xboxc2osc`

`xboxc2osc` is a tool for converting the input from an Xbox bluetooth controller to OSC signals

## Usage

When you run this application in the terminal, it will attempt to locate an Xbox Controller connected to your system via HID. If found, input from the controller will be redirected to OSC (see below),

Run `xboxc2osc -h` for available settings:

```
Usage of ./xboxc2osc:
  -host string
    	Osc reciever host (default "127.0.0.1")
  -pid string
    	HID product id as hexidecimal string (default "0b13")
  -port int
    	Osc reciever port (default 8000)
  -prefix string
    	Osc path prefix (default "xboxc")
  -vid string
    	HID vendor id as hexidecimal string (default "045e")
```

## OSC routes

By default the osc routes are prefixed with `/xboxc`, you can change this by setting the `-prefix` flag.

## Buttons

Button values are sent as a `true` when a press is detected, and reset to `false` after press.

 - `/xboxc/btn/a` _(bool)_
 - `/xboxc/btn/b` _(bool)_
 - `/xboxc/btn/x` _(bool)_
 - `/xboxc/btn/y` _(bool)_
 - `/xboxc/btn/lb` _(bool)_ the left index finger button (above left trigger)
 - `/xboxc/btn/rb` _(bool)_ the right index finger button (above the right trigger)
 - `/xboxc/btn/ls` _(bool)_ the left joystick button
 - `/xboxc/btn/rs` _(bool)_ the right joystick button
 - `/xboxc/btn/select` _(bool)_
 - `/xboxc/btn/start` _(bool)_


## Direction Pad

Input to the DPad is treated like button presses, sending `true` and then `false`. The route name depends on the compass direction (NW = Northwest, so the top left of the Dpad).

- `xboxc/btn/n` _(bool)_
- `xboxc/btn/ne` _(bool)_
- `xboxc/btn/e` _(bool)_
- `xboxc/btn/se` _(bool)_
- `xboxc/btn/s` _(bool)_
- `xboxc/btn/sw` _(bool)_
- `xboxc/btn/w` _(bool)_
- `xboxc/btn/nw` _(bool)_

## Joysticks

Joystick values are sent as floating point numbers between 0.0 and 1.0. There are composite routs that send an array of xy values, and the x and y components are broken into individual routes.

- `xboxc/btn/ls` left joystick
  - _(float)_ x value
  - _(float)_ y value
- `xboxc/btn/ls/x` _(float)_
- `xboxc/btn/ls/y` _(float)_
- `xboxc/btn/rs` right joystick
  - _(float)_ x value
  - _(float)_ y value
- `xboxc/btn/rs/x` _(float)_
- `xboxc/btn/rs/y` _(float)_

## Triggers

Trigger values are sent as 0.0 - 1.0 numbers, same as with the joysticks.

- `xboxc/btn/lt` _(float)_
- `xboxc/btn/rt` _(float)_







