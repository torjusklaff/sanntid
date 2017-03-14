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
	statesPort          = 16573
	globalQueuePort    = 16574
	broadcastTime       = 1 * time.Second
)

// Setter opp alle channels og funksjoner i en felles initialisering
func NetworkInit(
	id string,
	nElevators chan int,
	receiveNewOrder chan def.Order,
	receiveRemoverOrder chan def.Order,
	sendNewOrder chan def.Order,
	sendRemoveOrder chan def.Order,
	sendGlobalQueue chan [4][2]int,
	receivedGlobalQueue chan [4][2]int,
	sendStates chan def.Elevator,
	receivedStates chan def.Elevator) {

	go PeerListener(id, nElevators)
	go SendMsg(id, sendNewOrder, sendRemoveOrder, sendGlobalQueue, sendStates)
	go ReceiveMsg(id, receiveNewOrder, receiveRemoverOrder, receivedGlobalQueue, receivedStates)
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
func PeerListener(id string, nElevators chan int) {
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
			nElevators <- len(p.Peers)
			fmt.Printf("Number of active peers: %v \n", len(p.Peers))
		}
	}
}

// Setter opp channels for broadcast og sender det som kommer inn på input-channelsene
// se main fra network-module gitt på github
func SendMsg(
	localIP string,
	sendNewOrder chan def.Order,
	sendRemoveOrder chan def.Order,
	sendGlobalQueue chan [4][2]int,
	sendStates chan def.Elevator) {
	bcastSendNewOrder := make(chan def.Order)
	bcastSendRemoveOrder := make(chan def.Order)
	bcastSendGlobalQueue := make(chan [4][2]int)
	bcastSendStates := make(chan def.Elevator)

	go bcast.Transmitter(sendOrderPort, bcastSendNewOrder)
	go bcast.Transmitter(removeOrderPort, bcastSendRemoveOrder)
	go bcast.Transmitter(globalQueuePort, bcastSendGlobalQueue)
	go bcast.Transmitter(statesPort, bcastSendStates)
	
	for {
		select {
		case msg := <-sendNewOrder:
			sending := msg
			bcastSendNewOrder <- sending
		case msg := <-sendRemoveOrder:
			sending := msg
			bcastSendRemoveOrder <- sending
		case msg := <-sendGlobalQueue:
			sending := msg
			bcastSendGlobalQueue <- sending
		case msg := <-sendStates:
			sending := msg
			bcastSendStates <- sending
		default:
		}
	}
}

// Setter opp channels som lytter etter msg fra SendMsg()		(se main fra network-modul)
func ReceiveMsg(
	LocalIP string,
	receiveNewOrder chan def.Order,
	receiveRemoverOrder chan def.Order,
	receivedGlobalQueue chan [4][2]int,
	receivedStates chan def.Elevator) {

	bcastReceiveNewOrder := make(chan def.Order)
	bcastReceiveRemoveOrder := make(chan def.Order)
	bcastReceiveGlobalQueue := make(chan [4][2]int)
	bcastReceiveStates := make(chan def.Elevator)

	go bcast.Receiver(sendOrderPort, bcastReceiveNewOrder)
	go bcast.Receiver(removeOrderPort, bcastReceiveRemoveOrder)
	go bcast.Receiver(globalQueuePort, bcastReceiveGlobalQueue)
	go bcast.Receiver(statesPort, bcastReceiveStates)

	for {
		select {
		case msg := <-bcastReceiveNewOrder:
			receiveNewOrder <- msg
		case msg := <-bcastReceiveRemoveOrder:
			receiveRemoverOrder <- msg
		case msg := <-bcastReceiveGlobalQueue:
			receivedGlobalQueue <- msg
		case msg := <-bcastReceiveStates:
			if msg.Id != LocalIP{
			receivedStates <- msg
			}
		}
	}
}
