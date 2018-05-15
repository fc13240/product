package main

import (
	"helper/rpi"
)

func main() {
	tank := rpi.StartTank()
	tank.Speed(0)
	tank.Left()
	//tank.Reverse()

}
