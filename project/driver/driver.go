package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import (
	def "/definitions"
	"C"
)

func Set_motor_direction(dirn def.Motor_direction) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Set_button_lamp(button def.Button_type, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Get_button_signal(button def.Button_type, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Clear_all_lamps() {
	for floor := 0; floor < N_floors; floor++ {
		if floor < N_floors-1 {
			Set_button_lamp(def.Buttoncall_down, floor, 0)
		}
		if floor > 0 {
			Set_button_lamp(def.Buttoncall_up, floor, 0)
		}
		Set_button_lamp(def.Buttoncall_internal, floor, 0)
	}
}

func Elev_init() {
	C.elev_init()
	Clear_all_lamps()

	Set_motor_direction(def.Dir_down)
	for Get_floor_sensor_signal() == -1 {
	}
	Set_motor_direction(def.Dir_stop)
}
