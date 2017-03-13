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
	fmt.Printf("%v\n", driver.GetFloor_sensor_signal())


	button_pressed := make(chan def.Order_button)
	fmt.Printf("Made channel button_pressed\n")
	onFloor := make(chan int)
	fmt.Printf("Made channel onFloor\n")


	for {
		go driver.ElevatorOnFloor(onFloor, elevator)
		go driver.Check_all_buttons(button_pressed)
		
		select{
			case button_is_actually_pressed := <- button_pressed:
				fsm.FSMButtonPressed(button_is_actually_pressed, &elevator)
			
			case floor := <- onFloor:
				//fmt.Printf("Inside onFloor-channel case\n")
				fsm.FSMFloor_arrival(floor, &elevator)
			
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
	messageReceive := make(chan def.Network_message)
	go bcast.Transmitter(16569, message_transmit)
	go bcast.Receiver(16569, messageReceive)

	cost_transmit := make(chan def.cost_message)
	costReceive := make(chan def.cost_message)
	go bcast.Transmitter(16570, cost_transmit)
	go bcast.Receiver(16570, costReceive)

	newOrder_transmit := make(chan def.Order_button) 		// sjekke om vi trenger buffer
	newOrderReceive := make(chan def.Order_button)		// sjekke om vi trenger buffer
	go bcast.Transmitter(16571, newOrder_transmit)
	go bcast.Receiver(16571, newOrderReceive)


	button_pressed := make(chan def.Order_button)
	go driver.Check_all_buttons(button_pressed)



	
	go func(){
		// bestilling på denne heisen
		if button def.Order_button := <- button_pressed{
			newOrder_transmit <- button
			cost = fsm.FSMButtonPressed(button, elevator)

			var cost_msg def.Cost_message
			cost_msg.cost = cost
			cost_msg.id = elevator.id

			cost_transmit <- cost_msg
		}
		// bestilling på annen heis
		if newOrder := <- newOrderReceive{
			cost = fsm.FSMButtonPressed(button, elevator)

			var cost_msg def.Cost_message
			cost_msg.cost = cost
			cost_msg.id = elevator.id

			cost_transmit <- cost_msg
		}
	}


	go func() {
		costs_for_elevators = make(map[string]float32)
		if received_cost := <- costReceive{
			elem, ok = costs_for_elevators[received_cost.id]
			if not ok{
				costs_for_elevators[received_cost.id] = received_cost.cost
			}
		}

	}






	for {
		select {
		case msg := <- messageReceive:
			//det som skjer dersom vi har kontakt med omverdenen
		case <-time.After(5*time.Second):
			// det som skjer dersom vi ikke får noe inn på receive-channelen


	}
*/



}
