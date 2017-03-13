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
	peer_port            = 20100
	send_order_port      = 20012
	remove_order_port    = 16572
	sendCost_port       = 16573
	globalQueue_port    = 16574
	ElevatorStates_port = 16575
	broadcast_time       = 1 * time.Second
)

// Setter opp alle channels og funksjoner i en felles initialisering
func NetworkInit(
	id string,
	numElevators chan int,
	receiveCost chan def.Cost,
	receiveNewOrder chan def.Order,
	receive_remover_order chan def.Order,
	sendCost chan def.Cost,
	sendNewOrder chan def.Order,
	sendRemoveOrder chan def.Order,
	sendGlobalQueue chan [4][2]int,
	receivedGlobalQueue chan [4][2]int,
	receivedStates chan def.Elevator,
	sendStates chan def.Elevator) {

	go PeerListener(id, numElevators)
	go SendMsg(id, sendCost, sendNewOrder, sendRemoveOrder, sendGlobalQueue, sendStates)
	go ReceiveMsg(receiveCost, receiveNewOrder, receive_remover_order, receivedGlobalQueue, receivedStates)
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
	go peers.Transmitter(peer_port, id, peerTxEnable)
	go peers.Receiver(peer_port, peerUpdateCh)
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
	sendCost chan def.Cost,
	sendNewOrder chan def.Order,
	sendRemoveOrder chan def.Order,
	sendGlobalQueue chan [4][2]int,
	send_ElevatorStates chan def.Elevator) {

	bcast_sendCost := make(chan def.Cost)
	bcast_sendNewOrder := make(chan def.Order)
	bcast_sendRemoveOrder := make(chan def.Order)
	bcast_sendGlobalQueue := make(chan [4][2]int)
	bcast_sendStates := make(chan def.Elevator)

	go bcast.Transmitter(sendCost_port, bcast_sendCost)
	go bcast.Transmitter(send_order_port, bcast_sendNewOrder)
	go bcast.Transmitter(remove_order_port, bcast_sendRemoveOrder)
	go bcast.Transmitter(globalQueue_port, bcast_sendGlobalQueue)
	go bcast.Transmitter(ElevatorStates_port, bcast_sendStates)

	for {
		select {
		case msg := <-sendCost:
			sending := msg
			bcast_sendCost <- sending
		case msg := <-sendNewOrder:
			sending := msg
			bcast_sendNewOrder <- sending
		case msg := <-sendRemoveOrder:
			sending := msg
			bcast_sendRemoveOrder <- sending
		case msg := <-sendGlobalQueue:
			sending := msg
			bcast_sendGlobalQueue <- sending
		case msg := <-send_ElevatorStates:
			sending := msg
			bcast_sendStates <- sending
		}
	}
}

// Setter opp channels som lytter etter msg fra Send_msg()		(se main fra network-modul)
func ReceiveMsg(
	receiveCost chan def.Cost,
	receiveNewOrder chan def.Order,
	receive_remover_order chan def.Order,
	receivedGlobalQueue chan [4][2]int,
	received_ElevatorStates chan def.Elevator) {

	bcast_receiveCost := make(chan def.Cost)
	bcast_receiveNewOrder := make(chan def.Order)
	bcast_receiveRemoveOrder := make(chan def.Order)
	bcast_receive_globalQueue := make(chan [4][2]int)
	bcast_receive_states := make(chan def.Elevator)

	go bcast.Receiver(sendCost_port, bcast_receiveCost)
	go bcast.Receiver(send_order_port, bcast_receiveNewOrder)
	go bcast.Receiver(remove_order_port, bcast_receiveRemoveOrder)
	go bcast.Receiver(globalQueue_port, bcast_receive_globalQueue)
	go bcast.Receiver(ElevatorStates_port, bcast_receive_states)

	for {
		select {
		case msg := <-bcast_receiveCost:
			receiveCost <- msg
		case msg := <-bcast_receiveNewOrder:
			receiveNewOrder <- msg
		case msg := <-bcast_receiveRemoveOrder:
			receive_remover_order <- msg
		case msg := <-bcast_receive_globalQueue:
			receivedGlobalQueue <- msg
		case msg := <-bcast_receive_states:
			received_ElevatorStates <- msg
		}
	}
}
