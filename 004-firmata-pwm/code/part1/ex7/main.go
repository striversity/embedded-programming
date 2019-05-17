package main

// Circuit: esp8266 - 5 leds

import (
	"flag"
	"fmt"
	"math/rand"
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
	rn    = rand.New(rand.NewSource(time.Now().Unix()))
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	fmt.Printf("Using Firmata over %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	led0 := gpio.NewDirectPinDriver(board1, "5")
	led1 := gpio.NewDirectPinDriver(board1, "4")
	led2 := gpio.NewDirectPinDriver(board1, "0")
	led3 := gpio.NewDirectPinDriver(board1, "2")
	led4 := gpio.NewDirectPinDriver(board1, "14")

	work := func() {
		gobot.Every(1000*time.Millisecond, func() {
			led0.On()
			time.Sleep(500 * time.Millisecond)
			led0.Off()
		})

		gobot.Every(1000*time.Millisecond, func() {
			led1.On()
			time.Sleep(time.Duration(rn.Intn(1000)) * time.Millisecond)
			led1.Off()
		})

		gobot.Every(2*time.Second, func() {
			led2.PwmWrite(byte(rn.Intn(255)))
		})

		gobot.Every(50*time.Millisecond, func() {
			led3.PwmWrite(byte(rn.Intn(255)))
		})

		var b byte
		var coutnUp = true
		gobot.Every(6*time.Millisecond, func() {
			led4.PwmWrite(b)

			if coutnUp {
				b++
			} else {
				b--
			}

			if b == 255 || b == 0 {
				coutnUp = !coutnUp
			}
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{led0, led1, led2, led3, led4},
		work,
	)

	robot.Start()
}
