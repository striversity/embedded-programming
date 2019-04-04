package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/cu.usbserial-A700dXAd")
	led0 := gpio.NewLedDriver(firmataAdaptor, "2")
	led1 := gpio.NewLedDriver(firmataAdaptor, "4")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led0.On()
			time.Sleep(300*time.Millisecond)
			led0.Off()
		})

		gobot.Every(700*time.Millisecond, func() {
			led1.Toggle()
		})
	}

	robot := gobot.NewRobot("Multi-blinker",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led0, led1},
		work,
	)

	robot.Start()
}
