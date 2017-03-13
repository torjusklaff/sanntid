package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	"fmt"
	"time"
)

func FsmFloorArrival(newFloor int, elevator *def.Elevator) {
	if newFloor == -1 {
		fmt.Print("Run FSMFloor_arrival while not on floor\n")
	} else {
		//fmt.Print("FSMFloor_arrival\n")
		driver.SetFloorIndicator(newFloor)
		elevator.LastFloor = newFloor
		elevator.MotorStopTimer.Stop()
		switch elevator.ElevatorState {
		case def.Moving:
			if queue.ShouldStop(*elevator) {
				driver.SetMotorDirection(def.DirStop)
				queue.ClearAtFloor(elevator, newFloor)
				driver.ClearLightsAtFloor(elevator.LastFloor)
				driver.SetDoorOpenLamp(1)
				elevator.DoorTimer.Reset(3 * time.Second)
				fmt.Printf("Timer started\n")
				elevator.ElevatorState = def.StopOnFloor
			}
			break
		case def.Idle:
		default:
			break
		}
	}
}

func FsmNextOrder(elevator *def.Elevator, next_order def.Order) { //arbitrator decides where we should go next
	driver.SetButtonLamp(next_order, 1)

	switch elevator.ElevatorState {
	case def.Idle:
		queue.Enqueue(elevator, next_order)
		if next_order.Floor == elevator.LastFloor {
			queue.ClearAtFloor(elevator, elevator.LastFloor)
			driver.ClearLightsAtFloor(elevator.LastFloor)
			elevator.DoorTimer.Reset(3 * time.Second)
			driver.SetDoorOpenLamp(1)
			elevator.ElevatorState = def.StopOnFloor
		} else {
			if next_order.Floor > elevator.LastFloor {
				elevator.CurrentDirection = def.DirUp
				driver.SetMotorDirection(elevator.CurrentDirection)
			} else {
				elevator.CurrentDirection = def.DirDown
				driver.SetMotorDirection(elevator.CurrentDirection)
			}

		}
		if elevator.CurrentDirection == def.DirStop {
			elevator.ElevatorState = def.Idle
		} else {
			elevator.ElevatorState = def.Moving
			elevator.MotorStopTimer.Reset(4 * time.Second)
		}
	case def.Moving:
	case def.StopOnFloor:
		queue.ClearAtFloor(elevator, elevator.LastFloor)
		driver.ClearLightsAtFloor(elevator.LastFloor)
		elevator.DoorTimer.Reset(3 * time.Second)
	case def.MotorStop:
		if next_order.Type == def.ButtonInternal {
			queue.Enqueue(elevator, next_order)
		}

	default:
		break
	}
	queue.PrintQueue(*elevator)
}

func FsmOnDoorTimeout(elevator *def.Elevator) {
	queue.PrintQueue(*elevator)
	driver.SetDoorOpenLamp(0)
	switch elevator.ElevatorState {
	case def.StopOnFloor:
		elevator.CurrentDirection = queue.ChooseDirection(*elevator)
		driver.SetMotorDirection(elevator.CurrentDirection)

		if elevator.CurrentDirection == def.DirStop {
			elevator.ElevatorState = def.Idle
		} else {
			elevator.ElevatorState = def.Moving
			elevator.MotorStopTimer.Reset(8 * time.Second)
		}
		break
	case def.Idle:
		elevator.CurrentDirection = queue.ChooseDirection(*elevator)
		driver.SetMotorDirection(elevator.CurrentDirection)

		if elevator.CurrentDirection == def.DirStop {
			elevator.ElevatorState = def.Idle
		} else {
			elevator.ElevatorState = def.Moving
			elevator.MotorStopTimer.Reset(8 * time.Second)
		}
		break
	default:
		break
	}
}

func FsmWhereToNext(elevator def.Elevator) {
	switch elevator.ElevatorState {
	case def.StopOnFloor:
		elevator.CurrentDirection = queue.ChooseDirection(elevator)
		driver.SetMotorDirection(elevator.CurrentDirection)

		if elevator.CurrentDirection == def.DirStop {
			elevator.ElevatorState = def.Idle
		} else {
			elevator.ElevatorState = def.Moving
			elevator.MotorStopTimer.Reset(8 * time.Second)
		}
		break
	case def.Idle:
		elevator.CurrentDirection = queue.ChooseDirection(elevator)
		driver.SetMotorDirection(elevator.CurrentDirection)

		if elevator.CurrentDirection == def.DirStop {
			elevator.ElevatorState = def.Idle
		} else {
			elevator.ElevatorState = def.Moving
			elevator.MotorStopTimer.Reset(8 * time.Second)
		}
		break
	default:
		break
	}
}

func FsmMotorStop(elevator *def.Elevator) def.Elevator {
	elevator.CurrentDirection = def.DirStop
	driver.SetMotorDirection(def.DirStop)

	elev := driver.ElevInitFromBackup()
	return elev
}
