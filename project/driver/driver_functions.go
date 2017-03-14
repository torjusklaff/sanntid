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

func SetButtonLampFromInternalQueue(queue [4][3]int) {
	for f := 0; f < def.NFloors; f++ {
		for btn := 0; btn < def.NButtons; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.ButtonType(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func SetButtonLampFromGlobalQueue(queue [4][2]int) {
	for f := 0; f < def.NFloors; f++ {
		for btn := 0; btn < 2; btn++ {

			var button def.Order
			button.Floor = f
			button.Type = def.ButtonType(btn)

			SetButtonLamp(button, queue[f][btn])
		}
	}
}

func CheckAllButtons(externalButtonPressed chan def.Order, internalButtonPressed chan def.Order) {
	var pressedButton def.Order
	var buttonSignal def.Order
	for {
		for floor := 0; floor < def.NFloors; floor++ {
			for button := 0; button < def.NButtons; button++ {
				buttonSignal.Floor = floor
				buttonSignal.Type = def.ButtonType(button)

				if GetButtonSignal(buttonSignal) == 1 {
					pressedButton.Type = def.ButtonType(button)
					pressedButton.Floor = floor
					if pressedButton.Type == def.ButtoncallInternal {
						internalButtonPressed <- pressedButton
					} else {
						externalButtonPressed <- pressedButton
					}
					time.Sleep(50*time.Millisecond)
				}
			}
		}
	}
}

func ElevatorOnFloor(onFloor chan int, elevator def.Elevator) {
	for {
		if (GetFloorSensorSignal() != elevator.LastFloor) && (GetFloorSensorSignal() != -1) {

			onFloor <- GetFloorSensorSignal()
		}
	}
}

func ClearLightsAtFloor(floor int) {
	for btn := 0; btn < def.NButtons; btn++ {
		var button def.Order
		button.Type = def.ButtonType(btn)
		button.Floor = floor
		SetButtonLamp(button, 0)
	}
}

func ElevInit() def.Elevator {
	SetMotorDirection(def.DirStop)
	C.elevInit()
	//clearAllLamps()

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
	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	motorStopTimer := time.NewTimer(10 * time.Second)
	motorStopTimer.Stop()

	var elev def.Elevator
	elev.LastFloor = GetFloorSensorSignal()
	elev.CurrentDirection = def.DirStop
	elev.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elev.ElevatorState = def.Idle
	elev.DoorTimer = doorTimer
	elev.MotorStopTimer = motorStopTimer

	if _, err := os.Stat("log.txt"); err == nil {
		lastQueue := backup.ReadLastLine(24)
		elev.Queue = backup.QueueFromString(lastQueue)
	}

	return elev
}

