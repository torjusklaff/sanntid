package definitions

import (
	"os/exec"
	"time"
)

const (
	NFloors    = 4
	NButtons   = 3
	numElevators = 3
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
	Idle          ElevStates = iota
	StopOnFloor             //Not really necessary, look into it (Change to OnFloor)
	Moving
	MotorStop
)

type Elevator struct {
	LastFloor        int
	CurrentDirection MotorDirection
	Queue             [NFloors][NButtons]int
	ElevatorState    ElevStates
	Id                string
	DoorTimer        *time.Timer
	MotorStopTimer  *time.Timer
	CurrentOrder     Order
}

/*type ElevatorMsg struct {
	LastFloor        int
	CurrentDirection MotorDirection
	ElevatorState    ElevStates
	Id                string
}*/

type Cost struct {
	Cost          float64
	CurrentOrder Order
	Id            string
}

type Channels struct {
	numElevators chan int
	receiveNewOrder chan def.Order
	receiveRemoveOrder chan def.Order
	receivedGlobalQueue chan [4][2]int
	receivedStates chan def.Elevator
	sendNewOrder chan def.Order
	sendRemoveOrder chan def.Order
	assignedNewOrder chan def.Order
	sendGlobalQueue chan [4][2]int
	sendStates chan def.Elevator
	errorHandling chan string
}

var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "main.go")
