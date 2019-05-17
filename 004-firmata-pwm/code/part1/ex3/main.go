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

var period1 = 1000
var onDuration1 = period1 / 2

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	fmt.Printf("Using port %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	led0 := gpio.NewLedDriver(board1, "5")
	led1 := gpio.NewLedDriver(board1, "4")

	work := func() {
		go readCommand()

		gobot.Every(500*time.Millisecond, func() {
			led0.Toggle()
		})

		gobot.Every(time.Duration(period1)*time.Millisecond, func() {
			if 0 < onDuration1 {
				led1.On()
				time.Sleep(time.Duration(onDuration1) * time.Millisecond)
			}
			led1.Off()
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{led0, led1},
		work,
	)

	robot.Start()
}

func readCommand() {
	var percent uint // 0 - 100%
	// loop forever reading user input
	for {
		fmt.Print("Enter duty cycle percentage (0-100): ")
		fmt.Scanln(&percent)
		fmt.Println()

		percent %= 101 // force to 0 - 100
		onDuration1 = int(float32(period1) * (float32(percent) / 100.0))
		fmt.Printf("DEBUG - On for %v\n", onDuration1)
	}
}
