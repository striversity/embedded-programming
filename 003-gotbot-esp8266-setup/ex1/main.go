package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor("10.10.100.107:3030")
	led0 := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led0.Toggle()
		})
	}

	robot := gobot.NewRobot("esp8266-bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led0},
		work,
	)

	robot.Start()
}
