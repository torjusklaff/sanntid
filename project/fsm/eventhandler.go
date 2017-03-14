package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	"fmt"
	"time"
)

func EventHandler(elevator *def.Elevator, ch def.Channels){
	allExternalOrders := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}
	sendStatesTicker := time.NewTicker(100*time.Millisecond)

	onFloor := pollFloors()
	go SafeKill(ch.errorHandling)

	for {
		select {
		case floor := <-onFloor:
			FsmFloorArrival(floor, elevator)

		case <-elevator.DoorTimer.C:
			fmt.Printf("Timer stopped\n")
			FsmOnDoorTimeout(elevator)

		case newOrder := <-ch.receiveNewOrder:
			queue.AddOrderToGlobalQueue(ch.sendGlobalQueue, allExternalOrders, newOrder)

		case newOrder := <-ch.assignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				fmt.Print("Assigned new order\n")
				queue.Enqueue(elevator, newOrder)
				FsmNextOrder(elevator, newOrder)
			}
		case globalQueue := <-ch.receivedGlobalQueue:
			allExternalOrders = globalQueue

		case <-elevator.MotorStopTimer.C:
			fmt.Print("main: detected motorStop\n")
			errorMessage := "MOTORSTOP"
			ch.errorHandling <- errorMessage
			elevator.ElevatorState = def.MotorStop

		case err := <-ch.errorHandling:
			if err == "MOTORSTOP" {
				elevator = FsmMotorStop(elevator)

				var dummyOrder def.Order
				dummyOrder.Floor = 1
				dummyOrder.Type = def.ButtoncallInternal

				FsmNextOrder(elevator, dummyOrder)
			}
			if err == "PROGRAMCRASH" {
				def.Restart.Run()
			}
		case <- sendStatesTicker.C:
			sendStates <- elevator
		default:
			break
		}
	}
}