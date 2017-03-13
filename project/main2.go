package main

// Test-main for driver-files
import (
	//arb "./arbitrator"
	//"./driver"
	"./backup"
	def "./definitions"
	//"./fsm"
	//net "./network"
	q "./queue"
	"fmt"
	"time"
	"os"
	"os/signal"
	"log"
)

func main() {

	DoorTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()

	var elevator def.Elevator
	elevator.LastFloor = 1
	elevator.CurrentDirection = def.DirStop
	elevator.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elevator.ElevatorState = def.Idle

	var previousOrder def.Order
	previousOrder.Type = def.ButtonInternal
	previousOrder.Floor = elevator.LastFloor
	previousOrder.Id = "-"
	previousOrder.Internal = true


	// 		BACKUP KAN NÅ LAGRE TING I FIL, SAMT AT KØ-MODULEN KAN DECODE STRINGS TIL KØ-ARRAYS
	queueString := q.Queue_to_string(elevator)
	backup.To_backup(queueString)

	stringSize := len(queueString)
	last_line := backup.Read_last_line(int64(stringSize))
	fmt.Print(last_line)

	go func(){
		var c = make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Print("User terminated program.\n")
		var err = os.Remove("log.txt")
		if err != nil {
	        log.Fatalf("Error deleting file: %v", err)
	    }
		log.Fatal("User terminated program.\n")
	}()

	for{}
	//queue := q.Queue_from_string(last_line+"\n")
	


	// 	CHANNELS 
	/*
	numElevators := make(chan int)

	receiveCost := make(chan def.Cost)
	receiveNewOrder := make(chan def.Order)
	receiveRemoveOrder := make(chan def.Order)

	sendCost := make(chan def.Cost)
	sendNewOrder := make(chan def.Order)
	sendRemoveOrder := make(chan def.Order)
	assignedNewOrder := make(chan def.Order)

	button_pressed := make(chan def.Order)
	onFloor := make(chan int)

	id := net.Get_id()
	go net.Network_init(id, numElevators, receiveCost, receiveNewOrder, receiveRemoveOrder, sendCost, sendNewOrder, sendRemoveOrder)
	go arb.Arbitrator_init(elevator, id, receiveNewOrder, assignedNewOrder, receiveCost, sendCost, numElevators) // button_pressed må endres til receiveNewOrder

	go driver.Check_all_buttons(sendNewOrder)
	go driver.Elevator_onFloor(onFloor, elevator)


	for {
		select {
		case floor := <-onFloor:
			fsm.FSMFloor_arrival(floor, &elevator, DoorTimer)
		case <-DoorTimer.C:
			fmt.Printf("Timer stopped\n")
			fsm.FSM_on_door_timeout(&elevator)
		case newOrder := <-assignedNewOrder:
			fmt.Print("Assigned new order\n")
			queue.Enqueue(&elevator, newOrder)
			fsm.FSM_next_order(&elevator, newOrder, DoorTimer)
		default:
			break
		}
	}
	*/
}


