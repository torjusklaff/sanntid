package main

// Test-main for driver-files
import (
	//"./driver"
	//"./backup"
	def "./definitions"
	//"./fsm"
	net "./network"
	//"./queue"
	"time"
	"fmt"
)


func main(){
	// Åpne ny backup-fil
	// If file not found: lag ny fil, initialisering

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
		case cost := <- receive_cost:
			fmt.Printf("Cost: %v", cost.Cost)
		case order := <- receive_new_order:
			fmt.Printf("Order: %v", order.Floor)
		}
	}

}


func Testing_network_channels(send_cost chan def.Cost, send_new_order chan def.Order_button) {
	it := 1
	btn := def.Order_button{def.Buttoncall_down, 1, false, ""}
	//cost_msg := def.Cost{0, btn, ""}
	for {
		btn.Floor = it
		send_new_order <- btn
		time.Sleep(5*time.Second)
	}
}
