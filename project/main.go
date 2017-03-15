package main

import (
	arb "./arbitrator"
	def "./definitions"
	"./driver"
	"./fsm"
	net "./network"
)

func main() {
	elevator := driver.ElevatorInit()
	elevator.Id = net.GetId()

	channels := def.Channels{
		NumElevators:                make(chan int),
		ReceiveNewOrder:             make(chan def.Order),
		ReceivedFloorOrderCompleted: make(chan int),
		SendStates:                  make(chan def.ElevatorMsg),
		ReceivedStates:              make(chan def.ElevatorMsg),
		SendNewOrder:                make(chan def.Order),
		SendFloorOrderCompleted:     make(chan int),
		AssignedNewOrder:            make(chan def.Order),
		ErrorHandling:               make(chan string)}

	go net.NetworkInit(&elevator, channels)
	go arb.ArbitratorInit(elevator, channels)
	go fsm.EventHandler(&elevator, channels)
	for {

	}

}
