package main

// Circuit: arduino-and-motor-shield-rev3
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

var (
	maSpeed byte
	maInc   = 1
	counter int
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
	maDirGpio := gpio.NewLedDriver(board1, maDirPin)

	work := func() {
		// enable motor direction
		maDirGpio.On()

		gobot.Every(40*time.Millisecond, func() {
			log.Infof("Setting speed to %v\n", maSpeed)
			maSpeedGpio.PwmWrite(maSpeed)
			if maSpeed == 0 {
				maDirGpio.Toggle()
			}

			maSpeed = byte(int(maSpeed) + maInc)
			
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
		})
	}

	robot := gobot.NewRobot("my-robot",
		[]gobot.Connection{board1},
		[]gobot.Device{maSpeedGpio, maDirGpio},
		work,
	)

	robot.Start()
}
