package main

// Test-main for driver-files
import (
	"./driver"
	//"./backup"
	def "./definitions"
	"./fsm"
	//"./network"
	//"./queue"
	//"time"
	"fmt"
)


func main(){
	// Åpne ny backup-fil
	// If file not found: lag ny fil, initialisering

	elevator := driver.Elev_init()
	fmt.Printf("%v\n", driver.Get_floor_sensor_signal())


	button_pressed := make(chan def.Order_button)
	fmt.Printf("Made channel button_pressed\n")
	on_floor := make(chan int)
	fmt.Printf("Made channel on_floor\n")


	for {
		go driver.Elevator_on_floor(on_floor, elevator)
		go driver.Check_all_buttons(button_pressed)
		
		select{
			case button_is_actually_pressed := <- button_pressed:
				fsm.FSM_button_pressed(button_is_actually_pressed, &elevator)
			
			case floor := <- on_floor:
				//fmt.Printf("Inside on_floor-channel case\n")
				fsm.FSM_floor_arrival(floor, &elevator)
			
			default:
				break
		}
		//fmt.Printf("End of for-loop\n")
	}
	

	


/*	
	// Getting localIP
	id := localip.LocalIP()
	elevator.id = id

	// Channels for updating alive peers on network
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)



	// We make channels for sending and receiving our custom data types
message_transmit := make(chan def.Network_message)
	message_receive := make(chan def.Network_message)
	go bcast.Transmitter(16569, message_transmit)
	go bcast.Receiver(16569, message_receive)

	cost_transmit := make(chan def.cost_message)
	cost_receive := make(chan def.cost_message)
	go bcast.Transmitter(16570, cost_transmit)
	go bcast.Receiver(16570, cost_receive)

	new_order_transmit := make(chan def.Order_button) 		// sjekke om vi trenger buffer
	new_order_receive := make(chan def.Order_button)		// sjekke om vi trenger buffer
	go bcast.Transmitter(16571, new_order_transmit)
	go bcast.Receiver(16571, new_order_receive)


	button_pressed := make(chan def.Order_button)
	go driver.Check_all_buttons(button_pressed)



	
	go func(){
		// bestilling på denne heisen
		if button def.Order_button := <- button_pressed{
			new_order_transmit <- button
			cost = fsm.FSM_button_pressed(button, elevator)

			var cost_msg def.Cost_message
			cost_msg.cost = cost
			cost_msg.id = elevator.id

			cost_transmit <- cost_msg
		}
		// bestilling på annen heis
		if new_order := <- new_order_receive{
			cost = fsm.FSM_button_pressed(button, elevator)

			var cost_msg def.Cost_message
			cost_msg.cost = cost
			cost_msg.id = elevator.id

			cost_transmit <- cost_msg
		}
	}


	go func() {
		costs_for_elevators = make(map[string]float32)
		if received_cost := <- cost_receive{
			elem, ok = costs_for_elevators[received_cost.id]
			if not ok{
				costs_for_elevators[received_cost.id] = received_cost.cost
			}
		}

	}






	for {
		select {
		case msg := <- message_receive:
			//det som skjer dersom vi har kontakt med omverdenen
		case <-time.After(5*time.Second):
			// det som skjer dersom vi ikke får noe inn på receive-channelen


	}
*/



}
