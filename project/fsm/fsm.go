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
		fmt.Print("Run FSMFloorArrival while not on floor\n")
	} else {
		driver.SetFloorIndicator(newFloor)
		elevator.LastFloor = newFloor
		elevator.MotorStopTimer.Stop()
		switch elevator.ElevatorState {
		case def.Moving:
			if queue.ShouldStop(*elevator) {
				driver.SetMotorDirection(def.DirStop)
				queue.DeleteInternalQueuesAtFloor(elevator, newFloor)
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

func FsmNextOrder(elevator *def.Elevator, nextOrder def.Order) { //arbitrator decides where we should go next
	fmt.Print("FSMNextOrder\n")
	driver.SetButtonLamp(nextOrder, 1)

	switch elevator.ElevatorState {
	case def.Idle:
		queue.Enqueue(elevator, nextOrder)
		if nextOrder.Floor == elevator.LastFloor {
			queue.DeleteInternalQueuesAtFloor(elevator, elevator.LastFloor)
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
			fmt.Print("FSMNextOrder: Reset motorTimer\n")
		}
	case def.Moving:
		break
	case def.StopOnFloor:
		queue.DeleteInternalQueuesAtFloor(elevator, elevator.LastFloor)
		driver.ClearLightsAtFloor(elevator.LastFloor)
		elevator.DoorTimer.Reset(3 * time.Second)
	case def.MotorStop:
		if nextOrder.Type == def.ButtoncallInternal {
			queue.Enqueue(elevator, nextOrder)
		}

	default:
		break
	}
	queue.PrintQueue(*elevator)
}

func FsmOnDoorTimeout(elevator *def.Elevator) {
	fmt.Print("FSMOnDoorTimeout\n")
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
			fmt.Print("FSMWhereToNext: Reset motorTimer\n")
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
			fmt.Print("FSMOnDoorTimeout: Reset motorTimer\n")
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
			fmt.Print("FSMWhereToNext: Reset motorTimer\n")
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
			fmt.Print("FSMOnDoorTimeout: Reset motorTimer\n")
		}
		break
	default:
		break
	}
}

func FsmMotorStop(elevator *def.Elevator) def.Elevator {
	fmt.Print("FSMMotorStop\n")
	elevator.CurrentDirection = def.DirStop
	driver.SetMotorDirection(def.DirStop)

	elev := driver.ElevatorInitFromBackup()
	return elev
}

func pollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := driver.GetFloorSensorSignal()

		for {
			newFloor := driver.GetFloorSensorSignal()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}

func SafeKill(errorHandling chan string) {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	var err = os.Remove("log.txt")
	fmt.Print("User terminated program.\n\n")
	driver.SetMotorDirection(def.DirStop)

	for i := 0; i < def.NFloors; i++ {
		driver.ClearLightsAtFloor(i)
	}
	//def.Restart.Run()

	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	log.Fatal("\nUser terminated program.\n")

}
