// Blinking built-in led example for TinyGo
// Flash Arduino devices over serial
// cmd: tinygo flash -target=arduino -port /dev/tty.usbmodem1411201 tinygo/blink/main.go
package main

import (
    "github.com/tinygo-org/tinygo/src/machine"
    "time"
)

const(
	defaultPort = "/dev/tty.usbmodem1411201"
)

func main() {
    led := machine.Pin(3) // machine.LED
    led.Configure(machine.PinConfig{Mode: machine.PinOutput})
    for {
        led.Low()
        time.Sleep(time.Millisecond * 1000)

        led.High()
        time.Sleep(time.Millisecond * 1000)
    }
}