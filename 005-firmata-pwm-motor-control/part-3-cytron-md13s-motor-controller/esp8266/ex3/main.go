package main

// Circuit: esp8266-and-cytron-motor-controller
// Objective: motorA speed and direction control using DirectPinDriver
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

	board1 := firmata.NewTCPAdaptor(port)
	maSpeedGpio := gpio.NewDirectPinDriver(board1, maPWMPin)
	maDirGpio := gpio.NewLedDriver(board1, maDirPin)

	work := func() {
		// enable motor direction
		maDirGpio.On()

		gobot.Every(40*time.Millisecond, func() {
			maSpeed = byte(int(maSpeed) + maInc)
			log.Infof("Setting speed to %v\n", maSpeed)
			maSpeedGpio.PwmWrite(maSpeed)

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
				maDirGpio.Toggle()
			}
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{maSpeedGpio, maDirGpio},
		work,
	)

	robot.Start()
}
