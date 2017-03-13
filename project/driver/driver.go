package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"
import def "../definitions"
import "fmt"
import "time"
import "../backup"

func SetMotorDirection(dirn def.Motor_direction) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button def.Order, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button.Type), C.int(button.Floor), C.int(value))
}

func SetButtonLampFromInternalQueue(queue [4][3]int) {
	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.Button_type(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func SetButtonLampFromGlobalQueue(queue [4][2]int) {
	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < 2; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.Button_type(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func GetButtonSignal(button def.Order) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button.Type), C.int(button.Floor)))
}

func CheckAllButtons(external_button_pressed chan def.Order, internal_button_pressed chan def.Order) {
	var pressed_button def.Order
	var button_signal def.Order
	for {
		for floor := 0; floor < def.N_floors; floor++ {
			for button := 0; button < def.N_buttons; button++ {
				button_signal.Floor = floor
				button_signal.Type = def.Button_type(button)

				if GetButtonSignal(button_signal) == 1 {
					pressed_button.Type = def.Button_type(button)
					pressed_button.Floor = floor
					if pressed_button.Type == def.Buttoncall_internal {
						internal_button_pressed <- pressed_button
					} else {
						external_button_pressed <- pressed_button
					}
				}
			}
		}
	}
}

func GetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func ElevatorOnFloor(on_floor chan int, elevator def.Elevator) {
	for {
		if (GetFloorSensorSignal() != elevator.Last_floor) && (GetFloorSensorSignal() != -1) {

			on_floor <- GetFloorSensorSignal()
		}
	}
}

func ClearLightsAtFloor(floor int) {
	for btn := 0; btn < def.N_buttons; btn++ {
		var button def.Order
		button.Type = def.Button_type(btn)
		button.Floor = floor
		SetButtonLamp(button, 0)
	}
}

func ElevInit() def.Elevator {
	SetMotorDirection(def.Dir_stop)
	C.elev_init()
	//clear_all_lamps()

	SetMotorDirection(def.Dir_down)

	it := 0
	for GetFloorSensorSignal() == -1 {
		it += 1
		if it == 100000 {
			SetMotorDirection(def.Dir_up)
		}
	}
	fmt.Printf("Found floor in init\n")
	SetMotorDirection(def.Dir_stop)
	SetFloorIndicator(GetFloorSensorSignal())

	// Initializing an elevator-object
	door_timer := time.NewTimer(3 * time.Second)
	door_timer.Stop()
	motor_stop_timer := time.NewTimer(10 * time.Second)
	motor_stop_timer.Stop()

	var elev def.Elevator
	elev.Last_floor = GetFloorSensorSignal()
	elev.Current_direction = def.Dir_stop
	elev.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elev.Elevator_state = def.Idle
	elev.Door_timer = door_timer
	elev.Motor_stop_timer = motor_stop_timer

	return elev
}

func ElevInitFromBackup() def.Elevator {
	elev := ElevInit()

	last_queue := backup.ReadLastLine(24)
	fmt.Print(last_queue)
	elev.Queue = backup.QueueFromString(last_queue)
	return elev
}
