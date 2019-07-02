package main

// Circuit: arduino-and-l298n-motor-controller
// Objective: motorA speed using DirectPinDriver
//
// | Enable | Dir 1 | Dir 2 | Motor         |
// +--------+-------+-------+---------------+
// | 0      | X     | X     | Off           |
// | 1      | 0     | 0     | 0ff           |
// | 1      | 0     | 1     | On (forward)  |
// | 1      | 1     | 0     | On (backward) |
// | 1      | 1     | 1     | Off           |

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

URL: https://tronixlabs.com.au/news/tutorial-l298n-dual-motor-controller-module-2a-and-arduino/

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 10	   | Speed
A-Dir1        | DIR1 (Motor A) | 9	   | Direction
A-Dir2        | DIR2 (Motor A) | 8	   | Direction
B-Enable      | PWMA (Motor B) | 5	   | Speed
B-Dir1        | DIR1 (Motor B) | 7	   | Direction
B-Dir2        | DIR2 (Motor B) | 6	   | Direction

*/
const (
	maPWMPin  = "10"
	maDir1Pin = "9"
	maDir2Pin = "8"
	mbPWMPin  = "5"
	mbDir1Pin = "7"
	mbDir2Pin = "6"
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
	maDir1Gpio := gpio.NewDirectPinDriver(board1, maDir1Pin)
	maDir2Gpio := gpio.NewDirectPinDriver(board1, maDir2Pin)

	work := func() {
		// enable motor direction
		maDir1Gpio.On()
		maDir2Gpio.Off()

		gobot.Every(500*time.Millisecond, func() {
			maSpeed += 5
			log.Infof("Setting speed to %v\n", maSpeed)
			maSpeedGpio.PwmWrite(maSpeed)
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{maSpeedGpio, maDir1Gpio, maDir2Gpio},
		work,
	)

	robot.Start()
}
