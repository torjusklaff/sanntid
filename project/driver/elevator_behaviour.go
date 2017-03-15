package driver

import (
	def "../definitions"
	"../backup"
	"time"
)

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

func ElevatorOnFloor(onFloor chan int, elevator def.Elevator) {
	for {
		if (FloorSensorSignal() != elevator.LastFloor) && (FloorSensorSignal() != -1) {

			onFloor <- FloorSensorSignal()
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

func ElevatorInit() def.Elevator {
	SetMotorDirection(def.DirStop)
	C.ElevatorInit()
	SetMotorDirection(def.DirDown)

	for FloorSensorSignal() == -1 {
	}

	SetMotorDirection(def.DirStop)
	SetFloorIndicator(FloorSensorSignal())

	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()
	motorStopTimer := time.NewTimer(10 * time.Second)
	motorStopTimer.Stop()

	elev := def.Elevator{
		LastFloor: FloorSensorSignal(),
		CurrentDirection: def.DirStop,
		Queue: [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		ElevatorState: def.Idle,
		DoorTimer: doorTimer,
		MotorStopTimer: motorStopTimer
	}

	if _, err := os.Stat("log.txt"); err == nil {
		lastQueue := backup.ReadLastLine(24)
		elev.Queue = backup.QueueFromString(lastQueue)
		for floor := 0; floor < def.NumFloors; floor++ {
			for button := 0; button < def.NumButtons; button++ {
				if elev.Queue[floor][button] == 1 {
					setButton := def.Order{Type: def.ButtonType(button), Floor: floor}
					SetButtonLamp(setButton, 1)
				}
			}
		}
	}
	

	return elev
}

