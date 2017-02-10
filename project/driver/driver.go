package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"

const (
	N_floors  = int(C.N_FLOORS)
	N_buttons = int(C.N_BUTTONS)
)

type motor_direction int

const (
	dir_down motor_direction = -1
	dir_stop motor_direction = 0
	dir_up   motor_direction = 1
)

type Button_type int

const (
	Buttoncall_down     Button_type = 1
	Buttoncall_up       Button_type = 0
	Buttoncall_internal Button_type = 2
)

type order_button struct {
	Type  Button_type
	Floor int
}

func Set_motor_direction(dirn motor_direction) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Set_button_lamp(button Button_type, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Get_button_signal(button Button_type, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Clear_all_lamps() {
	for floor := 0; floor < N_floors; floor++ {
		if floor < N_floors-1 {
			Set_button_lamp(Buttoncall_down, floor, 0)
		}
		if floor > 0 {
			Set_button_lamp(Buttoncall_up, floor, 0)
		}
		Set_button_lamp(Buttoncall_internal, floor, 0)
	}
}

func Elev_init() {
	C.elev_init()
	Clear_all_lamps()

	Set_motor_direction(dir_down)
	for Get_floor_sensor_signal() == -1 {
	}
	Set_motor_direction(dir_stop)
}



