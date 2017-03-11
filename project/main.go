package main

// Test-main for driver-files
import (
	arb "./arbitrator"
	"./driver"
	//"./backup"
	def "./definitions"
	"./fsm"
	net "./network"
	"./queue"
	"fmt"
	"time"
	"os"
	"os/signal"
	"log"
)

func main() {
	all_external_orders := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}

	door_timer := time.NewTimer(3 * time.Second)
	door_timer.Stop()
	motor_stop_timer := time.NewTimer(10 * time.Second)
	motor_stop_timer.Stop()
	
	var elevator def.Elevator
	if _, err := os.Stat("log.txt"); err == nil {
  		elevator = driver.Elev_init_from_backup()
	} else {
		elevator = driver.Elev_init()
	}

	fmt.Printf("%v\n", driver.Get_floor_sensor_signal())

	var previous_order def.Order
	previous_order.Type = def.Buttoncall_internal
	previous_order.Floor = elevator.Last_floor

	// 	CHANNELS
	n_elevators := make(chan int)

	error_handling := make(chan string)

	receive_cost := make(chan def.Cost)
	receive_new_order := make(chan def.Order)
	receive_remove_order := make(chan def.Order)
	received_global_queue := make(chan [4][2]int)

	send_cost := make(chan def.Cost)
	send_new_order := make(chan def.Order)
	send_remove_order := make(chan def.Order)
	assigned_new_order := make(chan def.Order)
	send_global_queue := make(chan [4][2]int)

	//button_pressed := make(chan def.Order)
	on_floor := make(chan int)



	id := net.Get_id()
	go net.Network_init(id, n_elevators, receive_cost, receive_new_order, receive_remove_order, send_cost, send_new_order, send_remove_order, send_global_queue, received_global_queue)
	go arb.Arbitrator_init(elevator, id, receive_new_order, assigned_new_order, receive_cost, send_cost, n_elevators) // MÅ ENDRE ARBITRATOREN TIL Å OPPFØRE SEG ANNERLEDES

	go driver.Check_all_buttons(send_new_order)
	go driver.Elevator_on_floor(on_floor, elevator)
	go Safe_kill()

	for {
		select {
		case floor := <-on_floor:
			fsm.FSM_floor_arrival(floor, &elevator, door_timer, motor_stop_timer)
		
		case <-door_timer.C:
			fmt.Printf("Timer stopped\n")
			fsm.FSM_on_door_timeout(&elevator, motor_stop_timer)

		case new_order := <- receive_new_order:
			queue.Update_global_queue(send_global_queue, all_external_orders, new_order)
		
		case new_order := <-assigned_new_order:
			fmt.Print("Assigned new order\n")
			queue.Backup_internal_queue(elevator)
			queue.Enqueue(&elevator, new_order)
			fsm.FSM_next_order(&elevator, new_order, door_timer, motor_stop_timer)
			driver.Set_button_lamp_from_internal_queue(elevator.Queue)

		case global_queue := <-received_global_queue:
			all_external_orders = global_queue
			driver.Set_button_lamp_from_global_queue(all_external_orders)

		case <-motor_stop_timer.C:
			error_message := "MOTORSTOP"
			error_handling <- error_message
			elevator.Elevator_state = def.Motor_stop
			fsm.FSM_motor_stop(&elevator)
		
		default:
			break
		}
	}

	//Test for `en heis
	/*
		door_timer := time.NewTimer(3 * time.Second)
		door_timer.Stop()

		elevator := driver.Elev_init()

		button_pressed := make(chan def.Order)
		fmt.Printf("Made channel button_pressed\n")
		on_floor := make(chan int)
		fmt.Printf("Made channel on_floor\n")
		go driver.Check_all_buttons(button_pressed)
		go driver.Elevator_on_floor(on_floor, elevator)

		for {
			select {
			case floor := <-on_floor:
				fsm.FSM_floor_arrival(floor, &elevator, door_timer)
			case <-door_timer.C:
				fmt.Printf("Timer stopped\n")
				fsm.FSM_on_door_timeout(&elevator)
			default:
				break
			}
		}*/

	// Åpne ny backup-fil
	// If file not found: lag ny fil, initialisering

	/*													// Test av nettverk

	go Testing_network_channels(send_cost, send_new_order)

	for {
		select {
		case cost := <-receive_cost:
			fmt.Printf("Cost: %v \n", cost.Cost)
		case order := <-receive_new_order:
			fmt.Printf("Order: %v \n", order.Floor)
		}
	}*/

}

/*func Testing_network_channels(send_cost chan def.Cost, send_new_order chan def.Order) {
	it := 1
	btn := def.Order{def.Buttoncall_down, 1, false, ""}
	//cost_msg := def.Cost{0, btn, ""}
	for {
		btn.Floor = it
		send_new_order <- btn
		time.Sleep(2 * time.Second)
		it += 1
	}
}*/

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
	log.Fatal("User terminated program.\n")
}