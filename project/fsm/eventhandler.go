package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
)

func EventHandler(elevator *def.Elevator, ch def.Channels) {
	globalQueue := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	onFloor := pollFloors()
	go SafeKill(ch.ErrorHandling)

	for {
		select {
		case floor := <-onFloor:
			FsmFloorArrival(floor, elevator, ch.SendFloorOrderCompleted)

		case <-elevator.DoorTimer.C:
			FsmOnDoorTimeout(elevator)

		case floorToDelete := <-ch.ReceivedFloorOrderCompleted:
			queue.DeleteGlobalOrdersAtFloor(&globalQueue, floorToDelete)
			driver.ClearExternalButtonLampsAtFloor(floorToDelete)

		case newOrder := <-ch.ReceiveNewOrder:
			queue.UpdateGlobalQueue(&globalQueue, newOrder)
			driver.SetButtonLamp(newOrder, 1)

		case newOrder := <-ch.AssignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				queue.Enqueue(elevator, newOrder)
				FsmNextOrder(elevator, newOrder)
			}

		case <-elevator.MotorStopTimer.C:
			errorMessage := "MOTORSTOP"
			ch.ErrorHandling <- errorMessage
			elevator.ElevatorState = def.MotorStop

		case err := <-ch.ErrorHandling:
			if err == "MOTORSTOP" {
				*elevator = FsmMotorStop(elevator)

				var dummyOrder def.Order
				dummyOrder.Floor = 1
				dummyOrder.Type = def.ButtoncallInternal
				FsmNextOrder(elevator, dummyOrder)
			}
			if err == "PROGRAMCRASH" {
				def.Restart.Run()
			}
			if err == "DISCONNECTED" {
				driver.StopButton(1)
			}
			if err == "CONNECTED" {
				driver.StopButton(0)
			}

		default:
			break
		}
	}
}
