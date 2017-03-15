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
	peerPort                = 20100
	sendOrderPort           = 20012
	removeOrderPort         = 16572
	statesPort              = 16573
	globalCompleteOrderPort = 16574
	broadcastTime           = 1 * time.Second
)

func NetworkInit(elevator *def.Elevator, ch def.Channels) {
	go PeerListener(elevator.Id, ch.NumElevators)
	go SendMsg(elevator.Id, ch.SendNewOrder, ch.SendFloorOrderCompleted, ch.SendStates)
	go ReceiveMsg(elevator.Id, ch.ReceiveNewOrder, ch.ReceivedFloorOrderCompleted, ch.ReceivedStates)
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

func PeerListener(id string, NumElevators chan int) {
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
			NumElevators <- len(p.Peers)
			fmt.Printf("Number of active peers: %v \n", len(p.Peers))
		}
	}
}

func SendMsg(
	localIP string,
	SendNewOrder chan def.Order,
	SendFloorOrderCompleted chan int,
	SendStates chan def.ElevatorMsg) {

	bcastSendNewOrder := make(chan def.Order)
	bcastSendFloorOrderCompleted := make(chan int)
	bcastSendStates := make(chan def.ElevatorMsg)

	go bcast.Transmitter(sendOrderPort, bcastSendNewOrder)
	go bcast.Transmitter(globalCompleteOrderPort, bcastSendFloorOrderCompleted)
	go bcast.Transmitter(statesPort, bcastSendStates)

	for {
		select {
		case msg := <-SendNewOrder:
			sending := msg
			bcastSendNewOrder <- sending
		case msg := <-SendFloorOrderCompleted:
			sending := msg
			bcastSendFloorOrderCompleted <- sending
		case msg := <-SendStates:
			sending := msg
			bcastSendStates <- sending
		default:
		}
	}
}

func ReceiveMsg(
	LocalIP string,
	ReceiveNewOrder chan def.Order,
	ReceivedFloorOrderCompleted chan int,
	ReceivedStates chan def.ElevatorMsg) {

	bcastReceiveNewOrder := make(chan def.Order)
	bcastReceivedFloorOrderCompleted := make(chan int)
	bcastReceiveStates := make(chan def.ElevatorMsg)

	go bcast.Receiver(sendOrderPort, bcastReceiveNewOrder)
	go bcast.Receiver(globalCompleteOrderPort, bcastReceivedFloorOrderCompleted)
	go bcast.Receiver(statesPort, bcastReceiveStates)

	for {
		select {
		case msg := <-bcastReceiveNewOrder:
			ReceiveNewOrder <- msg
		case msg := <-bcastReceivedFloorOrderCompleted:
			ReceivedFloorOrderCompleted <- msg
		case msg := <-bcastReceiveStates:
			if msg.Id != LocalIP {
				ReceivedStates <- msg
			}
		}
	}
}
