package main

//Plass 15 ip: 148 plass 12 ip: 144 plass 2 ip: 149 plass 3 ip: 150
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
	elevator := driver.ElevatorInit()

	channels := def.Channels{
		numElevators: make(chan int)
		receiveNewOrder: make(chan def.Order)
		receiveRemoveOrder: make(chan def.Order)
		receivedGlobalQueue: make(chan [4][2]int)
		receivedStates: make(chan def.Elevator, 10)
		sendNewOrder: make(chan def.Order)
		sendRemoveOrder: make(chan def.Order)
		assignedNewOrder: make(chan def.Order)
		sendGlobalQueue: make(chan [4][2]int)
		sendStates: make(chan def.Elevator)
		errorHandling: make(chan string)
	}

	go net.NetworkInit(&elevator, channels)
	go arb.ArbitratorInit(elevator, channels)
	go driver.CheckAllButtons(channels)

	go fsm.EventHandler(&elevator, channels)

}

