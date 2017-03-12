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
	all_external_orders := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	test_timer := time.NewTimer(1 * time.Second)
	//test_timer.Stop()

	var elevator def.Elevator
	if _, err := os.Stat("log.txt"); err == nil {
		elevator = driver.ElevInitFromBackup()
		var dummy_order def.Order
		dummy_order.Floor = 1
		dummy_order.Type = def.Buttoncall_internal
		fsm.FsmNextOrder(&elevator, dummy_order)
	} else {
		elevator = driver.ElevInit()
	}

	fmt.Printf("%v\n", driver.GetFloorSensorSignal())

	var previous_order def.Order
	previous_order.Type = def.Buttoncall_internal
	previous_order.Floor = elevator.Last_floor

	// 	CHANNELS
	n_elevators := make(chan int)

	//error_handling := make(chan string)

	receive_cost := make(chan def.Cost)
	receive_new_order := make(chan def.Order)
	receive_remove_order := make(chan def.Order)
	received_global_queue := make(chan [4][2]int)
	received_states := make(chan def.Elevator)

	send_cost := make(chan def.Cost)
	send_new_order := make(chan def.Order)
	send_remove_order := make(chan def.Order)
	assigned_new_order := make(chan def.Order)
	send_global_queue := make(chan [4][2]int)
	send_states := make(chan def.Elevator)

	on_floor := pollFloors()
	error_handling := make(chan string)

	id := net.GetId()

	go net.NetworkInit(id, n_elevators, receive_cost, receive_new_order, receive_remove_order, send_cost, send_new_order, send_remove_order, send_global_queue, received_global_queue, send_states, received_states)
	go arb.ArbitratorInit(elevator, id, receive_new_order, assigned_new_order, send_states, received_states, n_elevators) // MÅ ENDRE ARBITRATOREN TIL Å OPPFØRE SEG ANNERLEDES

	go driver.CheckAllButtons(send_new_order, assigned_new_order)
	//go driver.Elevator_on_floor(on_floor, elevator)

	go SafeKill()

	test_it := 0
	floor_sense := 0
	for {
		test_it += 1
		if sensor := driver.GetFloorSensorSignal(); sensor != -1 {
			floor_sense = sensor
		}
		if test_it == 500000 {
			backup.BackupInternalQueue(elevator)
			driver.SetButtonLampFromInternalQueue(elevator.Queue)
			driver.SetButtonLampFromGlobalQueue(all_external_orders)
			test_it = 0
		}
		select {
		case floor := <-on_floor:
			fsm.FsmFloorArrival(floor, &elevator)

		case <-elevator.Door_timer.C:
			fmt.Printf("Timer stopped\n")
			//queue.ClearGlobalQueue(send_global_queue, all_external_orders, elevator.Last_floor)
			fsm.FsmOnDoorTimeout(&elevator)

		case new_order := <-receive_new_order:
			queue.UpdateGlobalQueue(send_global_queue, all_external_orders, new_order)

		case new_order := <-assigned_new_order:
			if elevator.Queue[new_order.Floor][int(new_order.Type)] == 0 {
				fmt.Print("Assigned new order\n")
				queue.Enqueue(&elevator, new_order)
				fsm.FsmNextOrder(&elevator, new_order)
			}
		case global_queue := <-received_global_queue:
			all_external_orders = global_queue

		case <-elevator.Motor_stop_timer.C:
			fmt.Print("main: detected motor_stop\n")
			error_message := "MOTORSTOP"
			error_handling <- error_message
			elevator.Elevator_state = def.Motor_stop

		case err := <-error_handling:
			if err == "MOTORSTOP" {
				elevator = fsm.FsmMotorStop(&elevator)

				var dummy_order def.Order
				dummy_order.Floor = 1
				dummy_order.Type = def.Buttoncall_internal

				fsm.FsmNextOrder(&elevator, dummy_order)
			}
			if err == "PROGRAM_CRASH" {
				def.Restart.Run()
			}
		case <-test_timer.C:

			fmt.Printf("Current floor: %v \t Floor sensor: %v\n", elevator.Last_floor, floor_sense)
			test_timer.Reset(1 * time.Second)
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
	driver.SetMotorDirection(def.Dir_stop)

	for i := 0; i < def.N_floors; i++ {
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
