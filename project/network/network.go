package network

import (
	"./bcast"
	"./localip"
	"./peers"
	"time"
	"../def"
)

const (
	peer_port = 15647
	get_order_port = 16571
	remove_order_port = 16572
	get_cost_port = 16573
	backup_port = 16574
	broadcast_time = 1*time.Second
)

type Cost_msg struct {
	Adsress string
	Data def.Cost
}

type Order_msg struct {
	Adsress string
	Data def.Order_button
}


// Setter opp alle channels og funksjoner i en felles initialisering
func Network_init(
	n_elevators chan<- int, 
	receive_cost chan<- def.Cost,
	receive_new_order chan<- def.Order_button,
	receive_remover_order chan<- def.Order_button,
	send_cost chan<- def.Cost,
	send_new_order chan<- def.Order_button,
	send_remove_order chan<- def.Order_button){

	id := Get_id()
	go Peer_listener(id, n_elevators)
	go Send_msg(id, send_cost, send_new_order,send_remove_order)
	go Receive_msg(receive_cost, receive_new_order,receive_remover_order)
}



func Get_id() id {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
}

// Setter opp en peer-listener som sjekker etter updates på levende heiser
func Peer_listener(id string, number_of_elevators chan<- int){
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
			number_of_peers <- len(p.Peers)
		}
	}
}

// Setter opp channels for broadcast og sender det som kommer inn på input-channelsene 
// se main fra network-module gitt på github
func Send_msg(
	localIP string, 
	send_cost <-chan def.Cost, 
	send_new_order <-chan def.Order_button,
	send_remove_order <-chan def.Order_button){

	bcast_send_cost := make(chan Cost_msg)
	bcast_send_new_order := make(chan Order_msg)
	bcast_send_remove_order := make(chan Order_msg)

	go bcast.Transmitter(get_cost_port, bcast_send_cost)
	go bcast.Transmitter(get_order_port, bcast_send_new_order)
	go bcast.Transmitter(remove_order_port, bcast_send_remove_order)

	for {
		select {
		case msg := <-send_cost:
			sending := Cost_msg{Address: localIP, Data: msg}
			bcast_send_cost <- sending
		case msg := <-send_new_order:
			sending := Order_msg{Address: localIP, Data: msg}
			bcast_send_cost <- sending
		case msg := <-send_remove_order:
			sending := Order_msg{Address: localIP, Data: msg}
			bcast_send_cost <- sending
		}
	}
}

// Setter opp channels som lytter etter msg fra Send_msg()		(se main fra network-modul)
func Receive_msg(
	receive_cost chan<- def.Cost, 
	receive_new_order chan<- def.Order_button,
	receive_remover_order chan<- def.Order_button){

	bcast_receive_cost := make(chan Cost_msg)
	bcast_receive_new_order := make(chan Order_msg)
	bcast_receive_remove_order := make(chan Order_msg)

	go bcast.Receiver(get_cost_port, bcast_receive_cost)
	go bcast.Receiver(get_order_port, bcast_receive_new_order)
	go bcast.Receiver(remove_order_port, bcast_receive_remove_order)

	for {
		select {
		case msg := <-bcast_receive_cost:
			receive_cost <_ msg.Data
		case msg := <-bcast_receive_new_order:
			receive_new_order <_ msg.Data
		case msg := <-bcast_receive_remover_order:
			receive_remover_order <_ msg.Data
		}
	}	
}