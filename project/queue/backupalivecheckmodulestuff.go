//viiii trenger: checkifalive? any more? Restart process
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type State struct {
	Tick int32
}

type Message struct {
	PrimaryState State
}

var NullState State = State{0}

func SendMessages(outgoing_message chan Message) {
	local, err := net.ResolveUDPAddr("udp", ":44556") // Change to 127.0.0.1 to work on laptop
	if err != nil {
		log.Fatal(err)
	}

	bcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:33445")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		log.Fatal(err)
	}

	for {
		m := <-outgoing_message
		buffer := &bytes.Buffer{}
		binary.Write(buffer, binary.BigEndian, m)

		_, err := conn.WriteToUDP(buffer.Bytes(), bcast)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Log(str string) {
	fmt.Println(str)
	log.Println(str)
}

func LaunchBackupProcess() {
	cmd := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Redirect log output
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Set up com channel with backup
	outgoing_message := make(chan Message)
	go SendMessages(outgoing_message)

	LaunchBackupProcess()

	// Launch master process
	Log("Launching master process")
	state := NullState

	// Perhaps it is a restart?
	if len(os.Args) > 1 {
		initial_state, _ := strconv.Atoi(os.Args[1])
		state.Tick = int32(initial_state)
		Log(fmt.Sprintf("MASTER restart @%d", state.Tick))
		Log(fmt.Sprintf("MASTER PRINT %d", state.Tick))
	}

	for {
		Log("MASTER starting work")
		time.Sleep(1 * time.Second)
		state.Tick++

		Log("MASTER updated")
		time.Sleep(1 * time.Second)
		outgoing_message <- Message{state}

		Log("MASTER sent state to BACKUP")
		time.Sleep(1 * time.Second)

		Log(fmt.Sprintf("MASTER PRINT STATE %d", state.Tick))
		Log("")
		time.Sleep(1 * time.Second)
	}





package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type State struct {
	Tick int32
}

type Message struct {
	PrimaryState State
}

var NullState State = State{0}

func ListenForMessages(incoming_message chan Message) {
	local, err := net.ResolveUDPAddr("udp", ":33445") // Change to 127.0.0.1 to work on laptop
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buffer := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		b := bytes.NewBuffer(buffer)
		m := Message{}
		binary.Read(b, binary.BigEndian, &m)
		incoming_message <- m
	}
}

func RestartMasterProcess(initial_state State) {
	arg := fmt.Sprintf("go run primary.go %d", initial_state.Tick)
	cmd := exec.Command("gnome-terminal", "-x", "sh", "-c", arg)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Redirect log output
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create com channel with master
	incoming_message := make(chan Message)
	go ListenForMessages(incoming_message)

	fmt.Println("Launching backup process")
	state := NullState

	// Wait for initial state
	select {
	case msg := <-incoming_message:
		state = msg.PrimaryState
		fmt.Println("BACKUP Received initial state update @", state.Tick)
	case <-time.After(7 * time.Second):
		fmt.Println("BACKUP Has primary not started?")
	}

	for {
		select {
		case <-time.After(7 * time.Second):
			fmt.Println("BACKUP Primary loss detected. Take over @", state.Tick)
			fmt.Println("BACKUP PRINT", state.Tick)
			RestartMasterProcess(state)
			return
		case msg := <-incoming_message:
			state = msg.PrimaryState
			fmt.Println("BACKUP Update received. Primary state @", state.Tick)
		}
	}
}
