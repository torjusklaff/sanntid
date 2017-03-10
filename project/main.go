package main

// Test-main for driver-files
import (
	"./driver"
	//"./backup"
	def "./definitions"
	"./fsm"
	//net "./network"
	//"./queue"
	"fmt"
	"time"
)

func main() {
	door_timer := time.NewTimer(3 * time.Second)
	door_timer.Stop()

	elevator := driver.Elev_init()
	fmt.Printf("%v\n", driver.Get_floor_sensor_signal())

	button_pressed := make(chan def.Order_button)
	fmt.Printf("Made channel button_pressed\n")
	on_floor := make(chan int)
	fmt.Printf("Made channel on_floor\n")
	//next_in_queue := make(chan def.Order_button)
	go driver.Check_all_buttons(button_pressed)
	go driver.Elevator_on_floor(on_floor, elevator)

	for {
		select {
		case button_is_actually_pressed := <-button_pressed:
			fsm.FSM_button_pressed(button_is_actually_pressed, &elevator)
		case floor := <-on_floor:
			fsm.FSM_floor_arrival(floor, &elevator, door_timer)
		case <-door_timer.C:
			fmt.Printf("Timer stopped\n")
			fsm.FSM_on_door_timeout(&elevator)
		default:
			break
		}
	}
	// Åpne ny backup-fil
	// If file not found: lag ny fil, initialisering

	/*
		n_elevators := make(chan int)

		receive_cost := make(chan def.Cost)
		receive_new_order := make(chan def.Order_button)
		receive_remove_order := make(chan def.Order_button)

		send_cost := make(chan def.Cost)
		send_new_order := make(chan def.Order_button)
		send_remove_order := make(chan def.Order_button)

		go net.Network_init(n_elevators, receive_cost, receive_new_order, receive_remove_order, send_cost, send_new_order, send_remove_order)
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

func Testing_network_channels(send_cost chan def.Cost, send_new_order chan def.Order_button) {
	it := 1
	btn := def.Order_button{def.Buttoncall_down, 1, false, ""}
	//cost_msg := def.Cost{0, btn, ""}
	for {
		btn.Floor = it
		send_new_order <- btn
		time.Sleep(2 * time.Second)
		it += 1
	}
}
