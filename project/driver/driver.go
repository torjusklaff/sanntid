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
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button def.Order, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button.Type), C.int(button.Floor), C.int(value))
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func ButtonSignal(button def.Order) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button.Type), C.int(button.Floor)))
}

func FloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func StopButton(value int){
	C.elev_set_stop_lamp(value)
}

