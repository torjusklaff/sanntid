package definitions

import (
	"os/exec"
	"time"
)

const (
	NumFloors        = 4
	NumButtons       = 3
	MotorStopTimeout = 4 * time.Second
	DoorTimeout      = 2 * time.Second
)

type MotorDirection int

const (
	DirDown MotorDirection = -1
	DirStop MotorDirection = 0
	DirUp   MotorDirection = 1
)

type ButtonType int

const (
	ButtoncallDown     ButtonType = 1
	ButtoncallUp       ButtonType = 0
	ButtoncallInternal ButtonType = 2
)

type Order struct {
	Type     ButtonType
	Floor    int
	Internal bool
	Id       string
}

func OrderToString(order Order) string {
	var intern string
	if order.Internal == true {
		intern = "true"
	} else {
		intern = "false"
	}
	return "Type: " + string(order.Type) + "  Floor: " + string(order.Floor) + "  Internal: " + intern + "  Id: " + order.Id
}

type ElevStates int

const (
	Idle ElevStates = iota
	StopOnFloor
	Moving
	MotorStop
	NotConnected
)

type Elevator struct {
	LastFloor        int
	CurrentDirection MotorDirection
	Queue            [NumFloors][NumButtons]int
	ElevatorState    ElevStates
	Id               string
	DoorTimer        *time.Timer
	MotorStopTimer   *time.Timer
	CurrentOrder     Order
}

type ElevatorMsg struct {
	LastFloor        int
	CurrentDirection MotorDirection
	ElevatorState    ElevStates
	Id               string
}

type Cost struct {
	Cost         float64
	CurrentOrder Order
	Id           string
}

type Channels struct {
	NumElevators                chan int
	ReceiveNewOrder             chan Order
	ReceivedFloorOrderCompleted chan int
	ReceivedStates              chan ElevatorMsg
	SendStates                  chan ElevatorMsg
	SendNewOrder                chan Order
	SendFloorOrderCompleted     chan int
	AssignedNewOrder            chan Order
	ErrorHandling               chan string
}

var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "main.go")
