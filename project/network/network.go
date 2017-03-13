package network

import (
	def "../definitions"
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	peerPort            = 20100
	sendOrderPort      = 20012
	removeOrderPort    = 16572
	sendCostPort       = 16573
	globalQueuePort    = 16574
	elevatorStatesPort = 16575
	broadcast_time       = 1 * time.Second
)

type NetworkChannels struct {
	numElevators 			chan int
	receiveCost 			chan def.Cost
	receiveNewOrder 		chan def.Order
	receiveRemoveOrder 		chan def.Order
	sendCost 				chan def.Cost
	sendNewOrder 			chan def.Order
	sendRemoveOrder 		chan def.Order
	sendGlobalQueue 		chan [4][2]int
	receivedGlobalQueue 	chan [4][2]int
	receivedStates 			chan def.Elevator
	sendStates 				chan def.Elevator
}

// Setter opp alle channels og funksjoner i en felles initialisering
func NetworkInit(
	id string,
	netChannels NetworkChannels) {

	go PeerListener(id, netChannels.numElevators)
	go SendMsg(id, netChannels)
	go ReceiveMsg(netChannels)
}

func GetId() string {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}

// Setter opp en peer-listener som sjekker etter updates på levende heiser
func PeerListener(id string, numElevators chan int) {
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(peerPort, id, peerTxEnable)
	go peers.Receiver(peerPort, peerUpdateCh)
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			numElevators <- len(p.Peers)
			fmt.Printf("Number of active peers: %v \n", len(p.Peers))
		}
	}
}

// Setter opp channels for broadcast og sender det som kommer inn på input-channelsene
// se main fra network-module gitt på github
func SendMsg(
	localIP string,
	nc NetworkChannels) {

	bcastSendCost := make(chan def.Cost)
	bcastSendNewOrder := make(chan def.Order)
	bcastSendRemoveOrder := make(chan def.Order)
	bcastSendGlobalQueue := make(chan [4][2]int)
	bcastSendStates := make(chan def.Elevator)

	go bcast.Transmitter(sendCostPort, bcastSendCost)
	go bcast.Transmitter(sendOrderPort, bcastSendNewOrder)
	go bcast.Transmitter(removeOrderPort, bcastSendRemoveOrder)
	go bcast.Transmitter(globalQueuePort, bcastSendGlobalQueue)
	go bcast.Transmitter(elevatorStatesPort, bcastSendStates)

	for {
		select {
		case msg := <-nc.sendCost:
			sending := msg
			bcastSendCost <- sending
		case msg := <-nc.sendNewOrder:
			sending := msg
			bcastSendNewOrder <- sending
		case msg := <-nc.sendRemoveOrder:
			sending := msg
			bcastSendRemoveOrder <- sending
		case msg := <-nc.sendGlobalQueue:
			sending := msg
			bcastSendGlobalQueue <- sending
		case msg := <-nc.sendElevatorStates:
			sending := msg
			bcastSendStates <- sending
		}
	}
}

// Setter opp channels som lytter etter msg fra Send_msg()		(se main fra network-modul)
func ReceiveMsg(nc NetworkChannels) {

	bcastReceiveCost := make(chan def.Cost)
	bcastReceiveNewOrder := make(chan def.Order)
	bcastReceiveRemoveOrder := make(chan def.Order)
	bcastReceiveGlobalQueue := make(chan [4][2]int)
	bcastReceiveStates := make(chan def.Elevator)

	go bcast.Receiver(sendCostPort, bcastReceiveCost)
	go bcast.Receiver(sendOrderPort, bcastReceiveNewOrder)
	go bcast.Receiver(removeOrderPort, bcastReceiveRemoveOrder)
	go bcast.Receiver(globalQueuePort, bcastReceiveGlobalQueue)
	go bcast.Receiver(elevatorStatesPort, bcastReceiveStates)

	for {
		select {
		case msg := <-bcastReceiveCost:
			nc.receiveCost <- msg
		case msg := <-bcastReceiveNewOrder:
			nc.receiveNewOrder <- msg
		case msg := <-bcastReceiveRemoveOrder:
			nc.receiveRemoveOrder <- msg
		case msg := <-bcastReceiveGlobalQueue:
			nc.receivedGlobalQueue <- msg
		case msg := <-bcastReceiveStates:
			nc.receivedStates <- msg
		}
	}
}
