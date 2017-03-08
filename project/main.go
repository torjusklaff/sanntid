// Test-main for driver-files
import (
	"/driver"
	"/backup"
	def "/definitions"
	"/fsm"
	"/network"
	"/queue"
)


func main(){
	elevator = driver.Elev_init()

	// Getting localIP
	id := localip.LocalIP()
	elevator.id = id

	// Channels for updating alive peers on network
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	message_transmit := make(chan def.Network_message)
	message_receive := make(chan def.Network_message)
	go bcast.Transmitter(16569, message_transmit)
	go bcast.Receiver(16569, message_receive)

	

	
}
