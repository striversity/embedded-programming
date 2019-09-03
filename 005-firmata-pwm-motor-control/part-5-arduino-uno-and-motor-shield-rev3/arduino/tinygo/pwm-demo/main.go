// PWM demo using built-in led example for TinyGo
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
    machine.InitPWM()
    led := machine.PWM{machine.Pin(3)}
    led.Configure()
    var level uint16
    for {
        level += 50
        led.Set(level)
        time.Sleep(time.Millisecond * 50)
    }
}

// cycleColor is just a placeholder until math/rand or some equivalent is working.
func cycleColor(color uint8) uint8 {
	if color < 10 {
		return color + 1
	} else if color < 200 {
		return color + 10
	} else {
		return 0
	}
}