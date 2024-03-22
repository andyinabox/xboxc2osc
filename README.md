# Xbox Controller Relay

A library for reading the input from an Xbox bluetooth controller and redirecting it to other places.

Existing tools are listed below, however you can implement the [`EventPublisher`](main.go) interface to send to other sources (for example: WebSockets).


## Tools

 - `xbox2osc` redirects all incoming controller signals to [OSC](https://opensoundcontrol.stanford.edu/index.html) ([see docs](cmd/xboxc2osc/README.md))


