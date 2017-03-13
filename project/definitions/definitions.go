package definitions

import (
	"os/exec"
	"time"
)

const (
	NumFloors    = 4
	NumButtons   = 3
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
	ButtonDown     ButtonType = 1
	ButtonUp       ButtonType = 0
	ButtonInternal ButtonType = 2
)

type Order struct {
	Type     ButtonType
	Floor    int
	Internal bool
	Id       string
}

func Order_to_string(order Order) string {
	var intern string
	if order.Internal == true {
		intern = "true"
	} else {
		intern = "false"
	}
	return "Type: " + string(order.Type) + "  Floor: " + string(order.Floor) + "  Internal: " + intern + "  Id: " + order.Id
}

type elevatorStates int

const (
	Idle          elevatorStates = iota
	StopOnFloor             //Not really necessary, look into it (Change to onFloor)
	Moving
	MotorStop
)

type Elevator struct {
	LastFloor 		int
	CurrentDirection MotorDirection
	Queue 			[NumFloors][NumButtons]int
	elevatorStates 	elevatorStates
	Id 				string
	DoorTimer 		*time.Timer
	MotorStopTimer  *time.Timer
	CurrentOrder 	Order
}


type Cost struct {
	Cost          float64
	CurrentOrder Order
	Id            string
}

var Restart = exec.Command("gnome-terminal", "-x", "sh", "-c", "main")
