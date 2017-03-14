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
	peerPort           = 20100
	sendOrderPort      = 20012
	removeOrderPort    = 16572
	statesPort         = 16573
	globalQueuePort    = 16574
	broadcastTime      = 1 * time.Second
)


func NetworkInit(elevator *def.Elevator, ch def.Channels) {

	var id string
	go func(){
		for{
			flag.StringVar(&id, "id", "", "id of this peer")
			flag.Parse()
			localIP, err := localip.LocalIP()
			if err != nil {
				fmt.Println(err)
				localIP = "DISCONNECTED"
				elevator.ElevatorState = def.NotConnected
				ch.ErrorHandling <- "DISCONNECTED"
			}
			if (err == nil) && (elevator.ElevatorState == def.NotConnected){
				elevator.ElevatorState == def.Idle
				ch.ErrorHandling <- "CONNECTED"
			}
			id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		}
	}()

	elevator.Id = id

	go PeerListener(id, ch.NumElevators)
	go SendMsg(id, ch.SendNewOrder, ch.SendRemoveOrder, ch.SendGlobalQueue, ch.SendStates)
	go ReceiveMsg(id, ch.ReceiveNewOrder, ch.receiveRemoverOrder, ch.ReceivedGlobalQueue, ch.ReceivedStates)
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
	SendRemoveOrder chan def.Order,
	SendGlobalQueue chan [4][2]int,
	SendStates chan def.Elevator) {


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
		case msg := <-SendNewOrder:
			sending := msg
			bcastSendNewOrder <- sending
		case msg := <-SendRemoveOrder:
			sending := msg
			bcastSendRemoveOrder <- sending
		case msg := <-SendGlobalQueue:
			sending := msg
			bcastSendGlobalQueue <- sending
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
	receiveRemoverOrder chan def.Order,
	ReceivedGlobalQueue chan [4][2]int,
	ReceivedStates chan def.Elevator) {


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
			ReceiveNewOrder <- msg
		case msg := <-bcastReceiveRemoveOrder:
			receiveRemoverOrder <- msg
		case msg := <-bcastReceiveGlobalQueue:
			ReceivedGlobalQueue <- msg
		case msg := <-bcastReceiveStates:
			if msg.Id != LocalIP{
			ReceivedStates <- msg
			}
		}
	}
}
