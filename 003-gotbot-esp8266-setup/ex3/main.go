package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewTCPAdaptor("10.10.100.107:3030")
	led0 := gpio.NewLedDriver(firmataAdaptor, "15")
	led1 := gpio.NewLedDriver(firmataAdaptor, "12")
	led2 := gpio.NewLedDriver(firmataAdaptor, "14")

	var br byte
	work := func() {
		gobot.Every(1*time.Second, func() {
			led0.Toggle()
		})
		gobot.Every(700*time.Millisecond, func() {
			led1.On()
			time.Sleep(300 * time.Millisecond)
			led1.Off()
		})
		gobot.Every(20*time.Millisecond, func() {
			br++
			led2.Brightness(br)
		})
	}

	robot := gobot.NewRobot("esp8266-bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led0, led1, led2},
		work,
	)

	robot.Start()
}
