package main

// Circuit: esp8266-and-l298n-motor-controller
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
	defaultPort = "10.10.100.175:3030"
)

/*

URL: https://tronixlabs.com.au/news/tutorial-l298n-dual-motor-controller-module-2a-and-arduino/

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 14	   | Speed
A-Dir1        | DIR1 (Motor A) | 5	   | Direction
A-Dir2        | DIR2 (Motor A) | 4	   | Direction
B-Enable      | PWMA (Motor B) | 12	   | Speed
B-Dir1        | DIR1 (Motor B) | 0	   | Direction
B-Dir2        | DIR2 (Motor B) | 2	   | Direction

*/
const (
	maPWMPin  = "14"
	maDir1Pin = "5"
	maDir2Pin = "4"
	mbPWMPin  = "12"
	mbDir1Pin = "0"
	mbDir2Pin = "2"
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
