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
	test_timer.Stop()

	var elevator def.Elevator
	if _, err := os.Stat("log.txt"); err == nil {
		elevator = driver.Elev_init_from_backup()
		var dummy_order def.Order
		dummy_order.Floor = 1
		dummy_order.Type = def.Buttoncall_internal
		fsm.FSM_next_order(&elevator, dummy_order)
	} else {
		elevator = driver.Elev_init()
	}

	fmt.Printf("%v\n", driver.Get_floor_sensor_signal())

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

	send_cost := make(chan def.Cost)
	send_new_order := make(chan def.Order)
	send_remove_order := make(chan def.Order)
	assigned_new_order := make(chan def.Order)
	send_global_queue := make(chan [4][2]int)

	on_floor := make(chan int)

	id := net.Get_id()

	go net.Network_init(id, n_elevators, receive_cost, receive_new_order, receive_remove_order, send_cost, send_new_order, send_remove_order, send_global_queue, received_global_queue)
	go arb.Arbitrator_init(elevator, id, receive_new_order, assigned_new_order, receive_cost, send_cost, n_elevators) // MÅ ENDRE ARBITRATOREN TIL Å OPPFØRE SEG ANNERLEDES

	go driver.Check_all_buttons(send_new_order)
	go driver.Elevator_on_floor(on_floor, elevator)
	go Safe_kill()

	test_it := 0

	for {
		test_it += 1
		if test_it == 500000 {
			backup.Backup_internal_queue(elevator)
			driver.Set_button_lamp_from_internal_queue(elevator.Queue)
			driver.Set_button_lamp_from_global_queue(all_external_orders)
			test_it = 0
		}

		select {
		case floor := <-on_floor:
			fsm.FSM_floor_arrival(floor, &elevator)

		case <-elevator.Door_timer.C:
			fmt.Printf("Timer stopped\n")
			queue.Clear_global_queue(send_global_queue, all_external_orders, elevator.Last_floor)
			fsm.FSM_on_door_timeout(&elevator)

		case new_order := <-receive_new_order:
			queue.Update_global_queue(send_global_queue, all_external_orders, new_order)

		case new_order := <-assigned_new_order:
			if elevator.Queue[new_order.Floor][int(new_order.Type)] == 0 {
				fmt.Print("Assigned new order\n")
				queue.Enqueue(&elevator, new_order)
				fsm.FSM_next_order(&elevator, new_order)
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
				elevator = fsm.FSM_motor_stop(&elevator)

				var dummy_order def.Order
				dummy_order.Floor = 1
				dummy_order.Type = def.Buttoncall_internal

				fsm.FSM_next_order(&elevator, dummy_order)
			}
			if err == "PROGRAM_CRASH" {
				def.Restart.Run()
			}

		default:
			break
		}
	}
}

func Safe_kill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	var err = os.Remove("log.txt")
	fmt.Print("User terminated program.\n\n")
	driver.Set_motor_direction(def.Dir_stop)
	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	log.Fatal("\nUser terminated program.\n")
}
