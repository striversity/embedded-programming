package main

// Circuit: esp8266 - 2 leds

import (
	"flag"
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	defaultPort = "10.10.100.100:3030"
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	fmt.Printf("Using port %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	led0 := gpio.NewLedDriver(board1, "5")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led0.Toggle()
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{led0},
		work,
	)

	robot.Start()
}
