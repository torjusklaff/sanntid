package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import (
	"C"
	def "../definitions"
)

func SetMotorDirection(dirn def.MotorDirection) {
	C.elevSetMotorDirection(C.elevMotorDirectionT(dirn))
}

func SetButtonLamp(button def.Order, value int) {
	C.elevSetButtonLamp(C.elevButtonTypeT(button.Type), C.int(button.Floor), C.int(value))
}

func SetFloorIndicator(floor int) {
	C.elevSetFloorIndicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
	C.elevSetDoorOpenLamp(C.int(value))
}

func ButtonSignal(button def.Order) int {
	return int(C.elevButtonSignal(C.elevButtonTypeT(button.Type), C.int(button.Floor)))
}

func FloorSensorSignal() int {
	return int(C.elevFloorSensorSignal())
}

