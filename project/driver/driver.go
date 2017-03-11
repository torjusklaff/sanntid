package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"
import def "../definitions"
import "fmt"
import "os"
import (
	"../backup"
	"../queue"
	"os/signal"
	"log"
)

func Set_motor_direction(dirn def.Motor_direction) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Set_button_lamp(button def.Order, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button.Type), C.int(button.Floor), C.int(value))
}

func Set_button_lamp_from_queue(queue [][]int, global_or_internal string) {
	if global_or_internal == "global"{
		n_buttons = 2
	} else if global_or_internal == "internal"{
		n_buttons = def.N_buttons
	} else {
		n_buttons = 2
		fmt.Print("Wrong use of Set_button_lamp_from_queue")
	}

	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < n_buttons; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.Button_type(btn)

			Set_button_lamp(button, queue[f][btn])
		}
	}
}

func Set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Get_button_signal(button def.Order) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button.Type), C.int(button.Floor)))
}

func Check_all_buttons(button_pressed chan def.Order) {
	var pressed_button def.Order
	var button_signal def.Order
	for {
		for floor := 0; floor < def.N_floors; floor++ {
			for button := 0; button < def.N_buttons; button++ {
				button_signal.Floor = floor
				button_signal.Type = def.Button_type(button)

				if Get_button_signal(button_signal) == 1 {
					pressed_button.Type = def.Button_type(button)
					pressed_button.Floor = floor

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

func Clear_lights_at_floor(floor int) {
	for btn := 0; btn < def.N_buttons; btn++ {
		var button def.Order
		button.Type = def.Button_type(btn)
		button.Floor = floor
		Set_button_lamp(button, 0)
	}
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

func Elev_init_from_backup() def.Elevator {
	elevator := Elev_init()

	last_queue := backup.Read_last_line(12)
	elevator.Queue = queue.Queue_from_string(last_queue)
	return elevator
}


func Safe_kill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	var err = os.Remove("log.txt")
	Set_motor_direction(def.Dir_stop)
	if err != nil {
        log.Fatalf("Error deleting file: %v", err)
    }
	log.Fatal("User terminated program.\n")
}
