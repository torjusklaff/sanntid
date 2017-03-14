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
	SendStatesTicker := time.NewTicker(100*time.Millisecond)

	onFloor := pollFloors()
	go SafeKill(ch.ErrorHandling)

	for {
		select {
		case floor := <-onFloor:
			FsmFloorArrival(floor, elevator, allExternalOrders, ch.SendGlobalQueue)

		case <-elevator.DoorTimer.C:
			FsmOnDoorTimeout(elevator)

		case newOrder := <-ch.ReceiveNewOrder:
			queue.AddOrderToGlobalQueue(ch.SendGlobalQueue, allExternalOrders, newOrder)

		case newOrder := <-ch.AssignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				queue.Enqueue(elevator, newOrder)
				FsmNextOrder(elevator, newOrder)
			}
		case globalQueue := <-ch.ReceivedGlobalQueue:
			allExternalOrders = globalQueue

		case <-elevator.MotorStopTimer.C:
			errorMessage := "MOTORSTOP"
			ch.ErrorHandling <- errorMessage
			elevator.ElevatorState = def.MotorStop

		case err := <-ch.ErrorHandling:
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
			if err == "DISCONNECTED"{
				driver.StopButton(1)
			}
			if err == "CONNECTED"{
				driver.StopButton(0)
			}

		case <- SendStatesTicker.C:
			SendStates <- elevator
		default:
			break
		}
	}
}