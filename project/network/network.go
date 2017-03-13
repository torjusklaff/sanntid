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
	send_cost_port       = 16573
	global_queue_port    = 16574
	elevator_states_port = 16575
	broadcast_time       = 1 * time.Second
)

// Setter opp alle channels og funksjoner i en felles initialisering
func NetworkInit(
	id string,
	n_elevators chan int,
	receive_cost chan def.Cost,
	receive_new_order chan def.Order,
	receive_remover_order chan def.Order,
	send_cost chan def.Cost,
	send_new_order chan def.Order,
	send_remove_order chan def.Order,
	send_global_queue chan [4][2]int,
	received_global_queue chan [4][2]int,
	received_states chan def.Elevator,
	send_states chan def.Elevator) {

	go PeerListener(id, n_elevators)
	go SendMsg(id, send_cost, send_new_order, send_remove_order, send_global_queue, send_states)
	go ReceiveMsg(receive_cost, receive_new_order, receive_remover_order, received_global_queue, received_states)
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
func PeerListener(id string, n_elevators chan int) {
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
			n_elevators <- len(p.Peers)
			fmt.Printf("Number of active peers: %v \n", len(p.Peers))
		}
	}
}

// Setter opp channels for broadcast og sender det som kommer inn på input-channelsene
// se main fra network-module gitt på github
func SendMsg(
	localIP string,
	send_cost chan def.Cost,
	send_new_order chan def.Order,
	send_remove_order chan def.Order,
	send_global_queue chan [4][2]int,
	send_elevator_states chan def.Elevator) {

	bcast_send_cost := make(chan def.Cost)
	bcast_send_new_order := make(chan def.Order)
	bcast_send_remove_order := make(chan def.Order)
	bcast_send_global_queue := make(chan [4][2]int)
	bcast_send_states := make(chan def.Elevator)

	go bcast.Transmitter(send_cost_port, bcast_send_cost)
	go bcast.Transmitter(send_order_port, bcast_send_new_order)
	go bcast.Transmitter(remove_order_port, bcast_send_remove_order)
	go bcast.Transmitter(global_queue_port, bcast_send_global_queue)
	go bcast.Transmitter(elevator_states_port, bcast_send_states)

	for {
		select {
		case msg := <-send_cost:
			sending := msg
			bcast_send_cost <- sending
		case msg := <-send_new_order:
			sending := msg
			bcast_send_new_order <- sending
		case msg := <-send_remove_order:
			sending := msg
			bcast_send_remove_order <- sending
		case msg := <-send_global_queue:
			sending := msg
			bcast_send_global_queue <- sending
		case msg := <-send_elevator_states:
			sending := msg
			bcast_send_states <- sending
		}
	}
}

// Setter opp channels som lytter etter msg fra Send_msg()		(se main fra network-modul)
func ReceiveMsg(
	receive_cost chan def.Cost,
	receive_new_order chan def.Order,
	receive_remover_order chan def.Order,
	received_global_queue chan [4][2]int,
	received_elevator_states chan def.Elevator) {

	bcast_receive_cost := make(chan def.Cost)
	bcast_receive_new_order := make(chan def.Order)
	bcast_receive_remove_order := make(chan def.Order)
	bcast_receive_global_queue := make(chan [4][2]int)
	bcast_receive_states := make(chan def.Elevator)

	go bcast.Receiver(send_cost_port, bcast_receive_cost)
	go bcast.Receiver(send_order_port, bcast_receive_new_order)
	go bcast.Receiver(remove_order_port, bcast_receive_remove_order)
	go bcast.Receiver(global_queue_port, bcast_receive_global_queue)
	go bcast.Receiver(elevator_states_port, bcast_receive_states)

	for {
		select {
		case msg := <-bcast_receive_cost:
			receive_cost <- msg
		case msg := <-bcast_receive_new_order:
			receive_new_order <- msg
		case msg := <-bcast_receive_remove_order:
			receive_remover_order <- msg
		case msg := <-bcast_receive_global_queue:
			received_global_queue <- msg
		case msg := <-bcast_receive_states:
			received_elevator_states <- msg
		}
	}
}
