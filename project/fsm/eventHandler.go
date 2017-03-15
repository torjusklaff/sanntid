package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	//"fmt"
)

func EventHandler(elevator *def.Elevator, ch def.Channels) {
	globalQueue := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	onFloor := pollFloors()
	go buttonChecker(ch)
	go safeKill(ch.ErrorHandling)

	for {
		select {
		case floor := <-onFloor:
			fsmFloorArrival(floor, elevator, ch.SendFloorOrderCompleted)

		case <-elevator.DoorTimer.C:
			fsmOnDoorTimeout(elevator)

		case floorToDelete := <-ch.ReceivedFloorOrderCompleted:
			queue.DeleteGlobalOrdersAtFloor(&globalQueue, floorToDelete)
			driver.ClearExternalButtonLampsAtFloor(floorToDelete)

		case newOrder := <-ch.ReceiveNewOrder:
			queue.UpdateGlobalQueue(&globalQueue, newOrder)
			driver.SetButtonLamp(newOrder, 1)

		case newOrder := <-ch.AssignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				queue.Enqueue(elevator, newOrder)
				fsmNextOrder(elevator, newOrder)
			}

		case <-elevator.MotorStopTimer.C:
			elevator.ElevatorState = def.MotorStop

		default:
			break
		}
	}
}
