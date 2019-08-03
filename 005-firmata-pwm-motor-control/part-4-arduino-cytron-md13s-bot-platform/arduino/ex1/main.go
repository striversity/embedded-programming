package main

// Circuit: arduino-and-cytron-motor-controller
// Objective: motorA speed using DirectPinDriver
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
	maSpeed byte
)

func main() {
	flag.Parse()

	port := defaultPort

	if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}

	log.Infof("Using port %v\n", port)

	board1 := firmata.NewAdaptor(port)
	maSpeedGpio := gpio.NewDirectPinDriver(board1, maPWMPin)
	maDirGpio := gpio.NewDirectPinDriver(board1, maDirPin)

	work := func() {
		// enable motor direction
		maDirGpio.On()

		gobot.Every(500*time.Millisecond, func() {
			maSpeed += 5
			log.Infof("Setting speed to %v\n", maSpeed)
			maSpeedGpio.PwmWrite(maSpeed)
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{maSpeedGpio, maDirGpio},
		work,
	)

	robot.Start()
}
