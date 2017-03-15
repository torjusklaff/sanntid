package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	"log"
	"os"
	"os/signal"
	"time"
)

func FsmFloorArrival(newFloor int, elevator *def.Elevator, SendFloorOrdersCompleted chan int) {
	if newFloor == -1 {
	} else {
		driver.SetFloorIndicator(newFloor)
		elevator.LastFloor = newFloor
		elevator.MotorStopTimer.Stop()
		switch elevator.ElevatorState {
		case def.Moving:
			if queue.ShouldStop(*elevator) {
				driver.SetMotorDirection(def.DirStop)
				driver.ClearLightsAtFloor(elevator.LastFloor)
				queue.DeleteInternalQueueAtFloor(elevator, newFloor)
				driver.SetDoorOpenLamp(1)
				elevator.DoorTimer.Reset(3 * time.Second)
				elevator.ElevatorState = def.StopOnFloor
				SendFloorOrdersCompleted <- newFloor
			}
			break
		case def.NotConnected:
			if queue.ShouldStop(*elevator) {
				driver.SetMotorDirection(def.DirStop)
				queue.DeleteInternalQueueAtFloor(elevator, newFloor)
				driver.ClearLightsAtFloor(elevator.LastFloor)
				driver.SetDoorOpenLamp(1)
				elevator.DoorTimer.Reset(3 * time.Second)
			}
			break
		case def.Idle:
		default:
			break
		}
	}
}

func FsmNextOrder(elevator *def.Elevator, nextOrder def.Order) {
	driver.SetButtonLamp(nextOrder, 1)

	switch elevator.ElevatorState {
	case def.Idle:
		queue.Enqueue(elevator, nextOrder)
		if nextOrder.Floor == elevator.LastFloor {
			queue.DeleteInternalQueueAtFloor(elevator, elevator.LastFloor)
			driver.ClearLightsAtFloor(elevator.LastFloor)
			elevator.DoorTimer.Reset(3 * time.Second)
			driver.SetDoorOpenLamp(1)
			elevator.ElevatorState = def.StopOnFloor
		} else {
			if nextOrder.Floor > elevator.LastFloor {
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
		break
	case def.StopOnFloor:
		queue.DeleteInternalQueueAtFloor(elevator, elevator.LastFloor)
		driver.ClearLightsAtFloor(elevator.LastFloor)
		elevator.DoorTimer.Reset(3 * time.Second)
	case def.MotorStop:
		if nextOrder.Type == def.ButtoncallInternal {
			queue.Enqueue(elevator, nextOrder)
		}
	case def.NotConnected:
		queue.Enqueue(elevator, nextOrder)
		if nextOrder.Floor == elevator.LastFloor {
			queue.DeleteInternalQueueAtFloor(elevator, elevator.LastFloor)
			driver.ClearLightsAtFloor(elevator.LastFloor)
			elevator.DoorTimer.Reset(3 * time.Second)
			driver.SetDoorOpenLamp(1)
		} else {
			if nextOrder.Floor > elevator.LastFloor {
				elevator.CurrentDirection = def.DirUp
				driver.SetMotorDirection(elevator.CurrentDirection)
				elevator.MotorStopTimer.Reset(4 * time.Second)
			} else {
				elevator.CurrentDirection = def.DirDown
				driver.SetMotorDirection(elevator.CurrentDirection)
				elevator.MotorStopTimer.Reset(4 * time.Second)
			}
		}
	default:
		break
	}
}

func FsmOnDoorTimeout(elevator *def.Elevator) {
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
		elevator.ElevatorState = def.StopOnFloor
		FsmOnDoorTimeout(elevator)
	case def.NotConnected:
		elevator.CurrentDirection = queue.ChooseDirection(*elevator)
		driver.SetMotorDirection(elevator.CurrentDirection)

		if !(elevator.CurrentDirection == def.DirStop) {
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

	elev := driver.ElevatorInit()
	return elev
}

func ButtonChecker(ch def.Channels) {
	var pressedButton def.Order
	var buttonSignal def.Order
	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for button := 0; button < def.NumButtons; button++ {
				buttonSignal.Floor = floor
				buttonSignal.Type = def.ButtonType(button)

				if driver.ButtonSignal(buttonSignal) == 1 {
					pressedButton.Type = def.ButtonType(button)
					pressedButton.Floor = floor
					if pressedButton.Type == def.ButtoncallInternal {
						ch.AssignedNewOrder <- pressedButton //internal orders avoid the arbitrator alltogether
					} else {
						ch.SendNewOrder <- pressedButton
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}
}

func pollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := driver.FloorSensorSignal()

		for {
			newFloor := driver.FloorSensorSignal()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}

func SafeKill(ErrorHandling chan string) {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	var err = os.Remove("log.txt")
	driver.SetMotorDirection(def.DirStop)

	for i := 0; i < def.NumFloors; i++ {
		driver.ClearLightsAtFloor(i)
	}

	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	log.Fatal("\nUser terminated program.\n")

}
