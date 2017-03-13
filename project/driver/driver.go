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

func SetMotorDirection(dirn def.MotorDirection) {
	C.elev_set_MotorDirection(C.elev_MotorDirection_t(dirn))
}

func SetButtonLamp(button def.Order, value int) {
	C.elev_set_button_lamp(C.elev_ButtonType_t(button.Type), C.int(button.Floor), C.int(value))
}

func SetButtonLampFromInternalQueue(queue [4][3]int) {
	for f := 0; f < def.NumFloors; f++ {
		for btn := 0; btn < def.NumButtons; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.ButtonType(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func SetButtonLampFromGlobalQueue(queue [4][2]int) {
	for f := 0; f < def.NumFloors; f++ {
		for btn := 0; btn < 2; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.ButtonType(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func SetFloorIndicator(floor int) {
	C.elev_setFloor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func GetButtonSignal(button def.Order) int {
	return int(C.elev_get_buttonSignal(C.elev_ButtonType_t(button.Type), C.int(button.Floor)))
}

func CheckAllButtons(externalButtonPressed chan def.Order, internalButtonPressed chan def.Order) {
	var pressedButton def.Order
	var buttonSignal def.Order
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for button := 0; button < def.NumButtons; button++ {
				buttonSignal.Floor = floor
				buttonSignal.Type = def.ButtonType(button)

				if GetButtonSignal(buttonSignal) == 1 {
					pressedButton.Type = def.ButtonType(button)
					pressedButton.Floor = floor
					if pressedButton.Type == def.ButtonInternal {
						internalButtonPressed <- pressedButton
					} else {
						externalButtonPressed <- pressedButton
					}
				}
			}
		}
	}
}

func GetFloorSensorSignal() int {
	return int(C.elev_getFloor_sensor_signal())
}

func ElevatorOnFloor(onFloor chan int, elevator def.Elevator) {
	for {
		if (GetFloorSensorSignal() != elevator.LastFloor) && (GetFloorSensorSignal() != -1) {

			onFloor <- GetFloorSensorSignal()
		}
	}
}

func ClearLightsAtFloor(floor int) {
	for btn := 0; btn < def.NumButtons; btn++ {
		var button def.Order
		button.Type = def.ButtonType(btn)
		button.Floor = floor
		SetButtonLamp(button, 0)
	}
}

func ElevInit() def.Elevator {
	SetMotorDirection(def.DirStop)
	C.elev_init()
	//clear_all_lamps()

	SetMotorDirection(def.DirDown)

	it := 0
	for GetFloorSensorSignal() == -1 {
		it += 1
		if it == 100000 {
			SetMotorDirection(def.DirUp)
		}
	}
	fmt.Printf("Found floor in init\n")
	SetMotorDirection(def.DirStop)
	SetFloorIndicator(GetFloorSensorSignal())

	// Initializing an elevator-object
	DoorTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()
	MotorStopTimer := time.NewTimer(10 * time.Second)
	MotorStopTimer.Stop()

	var elev def.Elevator
	elev.LastFloor = GetFloorSensorSignal()
	elev.CurrentDirection = def.DirStop
	elev.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elev.ElevatorState = def.Idle
	elev.DoorTimer = DoorTimer
	elev.MotorStopTimer = MotorStopTimer

	return elev
}

func ElevInitFromBackup() def.Elevator {
	elev := ElevInit()

	lastQueue := backup.ReadLastLine(24)
	fmt.Print(lastQueue)
	elev.Queue = backup.QueueFromString(lastQueue)
	return elev
}
