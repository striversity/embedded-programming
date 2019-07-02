package main

// Circuit: esp8266 - esp8266-l293d-motor-shield
// Objective: motorA speed using DirectPinDriver

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	defaultPort = "10.10.100.175:3030"
)

/*

URL: https://hackaday.io/project/8856-incubator-controller/log/29291-node-mcu-motor-shield

Motor Shield	NodeMCU 			DevKit		GPIO	Purpose
---------------+-------------------+-----------+-------+----------
D1 				PWMA (Motor A)		D1			5		Speed
D3 				DIRA (Motor A)		D3			0		Direction
D2 				PWMA (Motor B)		D2			4		Speed
D4 				DIRB (Motor B)		D4			2		Direction

*/
const (
	maPWMPin = "5"
	maDirPin = "0"
	mbPWMPin = "4"
	mbDirPin = "2"
)

var (
	maSpeed byte
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	log.Infof("Using port %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	maSpeedGpio := gpio.NewDirectPinDriver(board1, maPWMPin)

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			maSpeed += 5
			log.Infof("Setting speed to %v\n", maSpeed)
			maSpeedGpio.PwmWrite(maSpeed)
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{maSpeedGpio},
		work,
	)

	robot.Start()
}
