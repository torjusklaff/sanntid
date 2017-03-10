package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"
import def "../definitions"
import "time"
import "fmt"

func Set_motor_direction(dirn def.Motor_direction) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Set_button_lamp(button def.Order_button, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button.Type), C.int(button.Floor), C.int(value))
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Get_button_signal(button def.Order_button) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button.Type), C.int(button.Floor)))
}

func Check_all_buttons(button_pressed chan def.Order_button) {
	var pressed_button def.Order_button
	var button_signal def.Order_button
	for {
		for floor := 0; floor < def.N_floors; floor++ {
			for button := 0; button < def.N_buttons; button++ {
				button_signal.Floor = floor
				button_signal.Type = def.Button_type(button)

				if Get_button_signal(button_signal) == 1 {
					pressed_button.Type = def.Button_type(button)
					pressed_button.Floor = floor

					fmt.Printf("Found button\n")
					button_pressed <- pressed_button
				}
			}
		}
	}
}

func Get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Elevator_on_floor(on_floor chan int, elevator def.Elevator) {
	for {
		if (Get_floor_sensor_signal() != elevator.Last_floor) && (Get_floor_sensor_signal() != -1) {
			on_floor <- Get_floor_sensor_signal()
		}
	}
}

func Door_open_close() {
	Set_door_open_lamp(1)
	door_timer := time.NewTimer(3 * time.Second)
	<-door_timer.C
	Set_door_open_lamp(0)
}

func Elev_init() def.Elevator {
	Set_motor_direction(def.Dir_stop)
	C.elev_init()
	//clear_all_lamps()

	Set_motor_direction(def.Dir_down)
	for Get_floor_sensor_signal() == -1 {
	}
	fmt.Printf("Found floor in init\n")
	Set_motor_direction(def.Dir_stop)
	Set_floor_indicator(Get_floor_sensor_signal())

	// Initializing an elevator-object
	var elev def.Elevator
	elev.Last_floor = Get_floor_sensor_signal()
	elev.Current_direction = def.Dir_stop
	elev.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elev.Elevator_state = def.Idle

	return elev
}
