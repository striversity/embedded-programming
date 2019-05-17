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

var (
	level byte
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	fmt.Printf("Using Firmata over %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	led0 := gpio.NewLedDriver(board1, "5")
	led1 := gpio.NewDirectPinDriver(board1, "4")

	work := func() {
		go readCommand()

		gobot.Every(1*time.Second, func() {
			led0.Brightness(level)
			led1.PwmWrite(level)
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
	var percent int // 0 - 100%
	// loop forever reading user input
	for {
		fmt.Print("Enter duty cycle percentage (0-100): ")
		fmt.Scanln(&percent)
		fmt.Println()

		percent %= 101 // force to 0 - 100
		level = byte(float32(255) * (float32(percent) / 100.0))
		fmt.Printf("DEBUG - On for %v\n", level)
	}
}
