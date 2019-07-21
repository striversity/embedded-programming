package main

// Circuit: esp8266-and-cytron-motor-controller
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
	defaultPort = "10.10.100.175:3030"
)

/*

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Dir         | DIR  (Motor A) | 0	   | Direction
A-Speed       | PWMA (Motor A) | 14	   | Speed
B-Dir         | DIR1 (Motor B) | 13	   | Direction
B-Speed       | PWMA (Motor B) | 15	   | Speed

*/
const (
	maDirPin = "0"
	maPWMPin = "14"
	mbDirPin = "13"
	mbPWMPin = "15"
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
		motorA.Direction("forward")
		motorB.Direction("backward")
		
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
	// log.Infof("Setting %v speed to %v\n", m.Name(), motorSpeed[idx])
	m.Speed(motorSpeed[idx])

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
		if m.CurrentDirection == "forward" {
			m.Direction("backward")
		} else {
			m.Direction("forward")
		}
	}
}
