package main

// Circuit: arduino-and-lcytron-motor-controller
// Objective: motorA speed and direction control using MotorDriver
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

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Dir         | DIR  (Motor A) | 4	   | Direction
A-Speed       | PWMA (Motor A) | 5	   | Speed
B-Dir         | DIR1 (Motor B) | 2	   | Direction
B-Speed       | PWMA (Motor B) | 3	   | Speed

*/
const (
	maDirPin = "4"
	maPWMPin = "5"
	mbDirPin = "2"
	mbPWMPin = "3"
)

var (
	maSpeed    byte
	maInc      = 1
	counter    = 0
	dirCounter = 0
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

	work := func() {
		motorA.Direction("forward")

		gobot.Every(40*time.Millisecond, func() {
			maSpeed = byte(int(maSpeed) + maInc)
			log.Infof("Setting speed to %v\n", maSpeed)
			motorA.Speed(maSpeed)

			counter++
			if counter == 255 {
				counter = 0
				if maInc == 1 {
					maInc = 0
				} else if maInc == 0 {
					maInc = -1
				} else {
					maInc = 1
				}
			}
			dirCounter++
			if dirCounter == 765 {
				dirCounter = 0
				if motorA.CurrentDirection == "forward" {
					motorA.Direction("backward")
				} else {
					motorA.Direction("forward")
				}
			}
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{motorA},
		work,
	)

	robot.Start()
}
