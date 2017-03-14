package main

//Plass 15 ip: 148 plass 12 ip: 144 plass 2 ip: 149 plass 3 ip: 150
// Test-main for driver-files
import (
	arb "./arbitrator"
	"./backup"
	def "./definitions"
	"./driver"
	"./fsm"
	net "./network"
	"./queue"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	allExternalOrders := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	sendStatesTicker := time.NewTicker(100*time.Millisecond)

	elevator := driver.ElevatorInit()
	var dummyOrder def.Order
	dummyOrder.Floor = 1
	dummyOrder.Type = def.ButtoncallInternal
	fsm.FsmNextOrder(&elevator, dummyOrder)

	fmt.Printf("%v\n", driver.GetFloorSensorSignal())


	// 	CHANNELS
	channels := def.Channels{
		numElevators: make(chan int)
		receiveNewOrder: make(chan def.Order)
		receiveRemoveOrder: make(chan def.Order)
		receivedGlobalQueue: make(chan [4][2]int)
		receivedStates: make(chan def.Elevator, 10)
		sendNewOrder: make(chan def.Order)
		sendRemoveOrder: make(chan def.Order)
		assignedNewOrder: make(chan def.Order)
		sendGlobalQueue: make(chan [4][2]int)
		sendStates: make(chan def.Elevator)
		errorHandling: make(chan string)
	}
	
	onFloor := fsm.pollFloors()
	elevator.Id = net.GetId()

	go net.NetworkInit(elevator.Id, channels)
	go arb.ArbitratorInit(elevator, channels)
	go driver.CheckAllButtons(channels)
	go SafeKill()

	for {
		select {
		case floor := <-onFloor:
			fsm.FsmFloorArrival(floor, &elevator)

		case <-elevator.DoorTimer.C:
			fmt.Printf("Timer stopped\n")
			//queue.DeleteGlobalQueuesAtFloor(sendGlobalQueue, allExternalOrders, elevator.LastFloor)
			fsm.FsmOnDoorTimeout(&elevator)

		case newOrder := <-receiveNewOrder:
			queue.AddOrderToGlobalQueue(sendGlobalQueue, allExternalOrders, newOrder)

		case newOrder := <-assignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				fmt.Print("Assigned new order\n")
				queue.Enqueue(&elevator, newOrder)
				fsm.FsmNextOrder(&elevator, newOrder)
			}
		case globalQueue := <-receivedGlobalQueue:
			allExternalOrders = globalQueue

		case <-elevator.MotorStopTimer.C:
			fmt.Print("main: detected motorStop\n")
			errorMessage := "MOTORSTOP"
			errorHandling <- errorMessage
			elevator.ElevatorState = def.MotorStop

		case err := <-errorHandling:
			if err == "MOTORSTOP" {
				elevator = fsm.FsmMotorStop(&elevator)

				var dummyOrder def.Order
				dummyOrder.Floor = 1
				dummyOrder.Type = def.ButtoncallInternal

				fsm.FsmNextOrder(&elevator, dummyOrder)
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

func SafeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	var err = os.Remove("log.txt")
	fmt.Print("User terminated program.\n\n")
	driver.SetMotorDirection(def.DirStop)

	for i := 0; i < def.NFloors; i++ {
		driver.ClearLightsAtFloor(i)
	}
	def.Restart.Run()

	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	log.Fatal("\nUser terminated program.\n")

}

