package main

// Circuit: arduino-and-motor-shield-rev3
// Objective: dual speed and direction control using MotorDriver
//
// | PWM  | Dir | Motor         |
// +------+-----+---------------+
// | 0    | X   | Off           |
// | 1    | 0   | On (forward)  |
// | 1    | 1   | On (backward) |

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	defaultPort = "/dev/cu.usbmodem1411201"
)

/*
https://store.arduino.cc/usa/arduino-motor-shield-rev3
Motor Shield  | Arduino        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Dir         | DIR  (Motor A) | 12	   | Direction
A-Speed       | PWMA (Motor A) | 3	   | Speed
B-Dir         | DIR1 (Motor B) | 13	   | Direction
B-Speed       | PWMA (Motor B) | 11	   | Speed
*/
const (
	maDirPin = "12"
	maPWMPin = "3"
	mbDirPin = "13"
	mbPWMPin = "11"
)

const (
	maIndex = iota
	mbIndex
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

	log.Infof("Using port %v\n", port)

	board1 := firmata.NewAdaptor(port)
	motorA := gpio.NewMotorDriver(board1, maPWMPin)
	motorA.DirectionPin = maDirPin
	motorA.SetName("Motor-A")
	motorB := gpio.NewMotorDriver(board1, mbPWMPin)
	motorB.DirectionPin = mbDirPin
	motorB.SetName("Motor-B")

	motors[maIndex] = motorA
	motors[mbIndex] = motorB

	work := func() {
		motorA.Off()
		motorB.Off()
		motorA.Direction("forward")
		motorB.Direction("forward")
		motorA.On()
		motorB.On()

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

	// log.Infof("Setting %v speed to %v\n", m.Name(), motorSpeed[idx])
	m.Speed(motorSpeed[idx])

	if motorSpeed[idx] == 0 {
		if m.CurrentDirection == "forward" {
			m.Direction("backward")
		} else {
			m.Direction("forward")
		}
	}

	motorSpeed[idx] = byte(int(motorSpeed[idx]) + motorInc[idx])

	counter[idx]++
	if counter[idx] == 255 {
		counter[idx] = 0
		if motorInc[idx] == 1 {
			motorInc[idx] = 0
		} else if motorInc[idx] == 0 {
			motorInc[idx] = -1
		} else {
			motorInc[idx] = 1
		}
	}
}
