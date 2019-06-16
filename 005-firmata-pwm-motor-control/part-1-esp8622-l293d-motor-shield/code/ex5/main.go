package main

// Circuit: esp8266 - esp8266-l293d-motor-shield
// Objective: dual speed and direction control using MotorDriver

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
const (
	maIndex = iota
	mbIndex = iota
)

var (
	motorSpeed [2]byte
	motorInc   = [2]int{1, 1}
	counter    = [2]int{}
	motors     [2]*gpio.MotorDriver
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	log.Info("Using port %v\n", port)

	board1 := firmata.NewTCPAdaptor(port)
	motorA := gpio.NewMotorDriver(board1, maPWMPin)
	motorA.DirectionPin = maDirPin
	motorA.SetName("Motor-A")
	motorB := gpio.NewMotorDriver(board1, mbPWMPin)
	motorB.DirectionPin = mbDirPin
	motorB.SetName("Motor-B")
	
	motors[maIndex] = motorA
	motors[mbIndex] = motorB

	work := func() {
		gobot.Every(40*time.Millisecond, func() {
			motorControl(maIndex)
		})

		gobot.Every(20*time.Millisecond, func() {
			motorControl(mbIndex)
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{motorA, motorB},
		work,
	)

	robot.Start()
}

func motorControl(idx int) {
	m := motors[idx]

	motorSpeed[idx] = byte(int(motorSpeed[idx]) + motorInc[idx])
	log.Info(m.Name(), " speed to ", motorSpeed[idx])
	motors[idx].Speed(motorSpeed[idx])

	counter[idx]++
	if counter[idx]%256 == 255 {
		if motorInc[idx] == 1 {
			motorInc[idx] = 0
		} else if motorInc[idx] == 0 {
			motorInc[idx] = -1
		} else {
			motorInc[idx] = 1
		}
	}

	if counter[idx]%766 == 765 {
		d := motors[idx].CurrentDirection
		if d == "forward" {
			d = "backward"
		} else {
			d = "forward"
		}
		motors[idx].Direction(d)
	}
}
