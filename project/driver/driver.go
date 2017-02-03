package Driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"


type motor_direction int

const (
	dir_down motor_direction = -1
	dir_stop motor_direction = 0
	dir_up motor_direction = 1
)


type button_type int

const (
	buttoncall_down button_type = 1
	buttoncall_up button_type = 0
	buttoncall_internal button_type = 2
)


func set_motor_direction(dirn motor_direction){
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func set_button_lamp(button button_type, floor int, value int){
	C.elev_set_button_lamp(C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value)))
}

func set_floor_indicator(floor int){
	C.elev_set_floor_indicator(C.int(floor))
}

func set_door_open(value int){
	C.elev_set_door_open(C.int(value))
}

func get_button_signal(button button_type, floor int){
	C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor))
}

func get_floor_sensor_signal(){
	C.elev_get_floor_sensor_signal()
}






