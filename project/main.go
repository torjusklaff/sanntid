package main

//Plass 15 ip: 148 plass 12 ip: 144
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

	testTimer := time.NewTimer(1 * time.Second)
	//testTimer.Stop()
	sendStatesTicker := time.NewTicker(100*time.Millisecond)

	var elevator def.Elevator
	if _, err := os.Stat("log.txt"); err == nil {
		elevator = driver.ElevInitFromBackup()
		var dummyOrder def.Order
		dummyOrder.Floor = 1
		dummyOrder.Type = def.ButtonInternal
		fsm.FsmNextOrder(&elevator, dummyOrder)
	} else {
		elevator = driver.ElevInit()
	}



	// 	CHANNELS
	netChannels := net.NetworkChannels{
		numElevators: make(chan int),
		receiveCost: make(chan def.Cost),
		receiveNewOrder: make(chan def.Order),
		receiveRemoveOrder: make(chan def.Order),
		receivedGlobalQueue: make(chan [4][2]int),
		receivedStates: make(chan def.Elevator),
		sendCost: make(chan def.Cost),
		sendNewOrder: make(chan def.Order),
		sendRemoveOrder: make(chan def.Order),
		assignedNewOrder: make(chan def.Order),
		sendGlobalQueue: make(chan [4][2]int),
		sendStates: make(chan def.Elevator)
	}

	onFloor := pollFloors()
	errorHandling := make(chan string)

	id := net.GetId()

	go net.NetworkInit(id, netChannels)
	go arb.ArbitratorInit(elevator, id, netChannels) // MÅ ENDRE ARBITRATOREN TIL Å OPPFØRE SEG ANNERLEDES

	go driver.CheckAllButtons(netChannels)

	go SafeKill()

	testIt := 0
	floorSense := 0
	for {
		testIt += 1
		if sensor := driver.GetFloorSensorSignal(); sensor != -1 {
			floorSense = sensor
		}
		if testIt == 500000 {
			backup.BackupInternalQueue(elevator)
			driver.SetButtonLampFromInternalQueue(elevator.Queue)
			driver.SetButtonLampFromGlobalQueue(allExternalOrders)
			testIt = 0
		}
		select {
		case floor := <-onFloor:
			fsm.FsmFloorArrival(floor, &elevator)
			sendStates <- elevator

		case <-elevator.DoorTimer.C:
			fsm.FsmOnDoorTimeout(&elevator)

		case newOrder := <-receiveNewOrder:
			queue.UpdateGlobalQueue(netChannels.sendGlobalQueue, allExternalOrders, newOrder)

		case newOrder := <-assignedNewOrder:
			if elevator.Queue[newOrder.Floor][int(newOrder.Type)] == 0 {
				queue.Enqueue(&elevator, newOrder)
				fsm.FsmNextOrder(&elevator, newOrder)
			}
		case globalQueue := <-receivedGlobalQueue:
			allExternalOrders = globalQueue

		case <-elevator.MotorStopTimer.C:
			error_message := "MOTORSTOP"
			errorHandling <- error_message
			elevator.ElevatorState = def.MotorStop

		case err := <-errorHandling:
			if err == "MOTORSTOP" {
				elevator = fsm.FsmMotorStop(&elevator)

				var dummyOrder def.Order
				dummyOrder.Floor = 1
				dummyOrder.Type = def.ButtonInternal

				fsm.FsmNextOrder(&elevator, dummyOrder)
			}
			if err == "PROGRAM_CRASH" {
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
	driver.SetMotorDirection(def.DirStop)

	for i := 0; i < def.NumFloors; i++ {
		driver.ClearLightsAtFloor(i)
	}

	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	log.Fatal("\nUser terminated program.\n")

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
