package main

import (
	arb "./arbitrator"
	def "./definitions"
	net "./network"
	"./fsm"
	"./driver"
)

func main() {
	elevator := driver.ElevatorInit()

	channels := def.Channels{
		NumElevators: make(chan int),
		ReceiveNewOrder: make(chan def.Order),
		ReceiveRemoveOrder: make(chan def.Order),
		ReceivedGlobalQueue: make(chan [4][2]int),
		ReceivedStates: make(chan def.Elevator, 10),
		SendNewOrder: make(chan def.Order),
		SendRemoveOrder: make(chan def.Order),
		SendGlobalQueue: make(chan [4][2]int),
		AssignedNewOrder: make(chan def.Order),
		SendStates: make(chan def.Elevator),
		ErrorHandling: make(chan string)
	}

	go net.NetworkInit(&elevator, channels)
	go arb.ArbitratorInit(elevator, channels)
	go fsm.ButtonChecker(channels)
	go fsm.EventHandler(&elevator, channels)

}

