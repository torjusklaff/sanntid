package queue

import (
	"../backup"
	def "../definitions"
	"fmt"
)

func requestsAbove(e def.Elevator) bool {
	for f := e.LastFloor + 1; f < def.NumFloors; f++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if e.Queue[f][def.ButtonType(btn)] == 1 {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e def.Elevator) bool {
	for f := 0; f < e.LastFloor; f++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if e.Queue[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func ChooseDirection(e def.Elevator) def.MotorDirection {
	switch e.CurrentDirection {
	case def.DirUp:
		if requestsAbove(e) {
			return def.DirUp
		} else if requestsBelow(e) {
			return def.DirDown
		} else {
			return def.DirStop
		}
	case def.DirDown:
		if requestsBelow(e) {
			return def.DirDown
		} else if requestsAbove(e) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	case def.DirStop:
		if requestsBelow(e) {
			return def.DirDown
		} else if requestsAbove(e) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	default:
		return def.DirStop
	}
	return def.DirStop
}

func ClearAtFloor(e *def.Elevator, floor int) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if e.Queue[floor][btn] == 1 {
			e.Queue[floor][btn] = 0
			backup.BackupInternalQueue(*e)
		}
	}
}

func ClearGlobalQueue(sendGlobalQueue chan [4][2]int, old_queue [4][2]int, floor int) {
	for btn := 0; btn < 2; btn++ {
		if old_queue[floor][btn] == 1 {
			old_queue[floor][btn] = 0
		}
	}
	sendGlobalQueue <- old_queue
}

func PrintQueue(e def.Elevator) {
	for f := 0; f < def.NumFloors; f++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			fmt.Printf("%v ", e.Queue[f][btn])
		}
		fmt.Printf("\n")
	}
}

func ShouldStop(e def.Elevator) bool {
	switch e.CurrentDirection {
	case def.DirDown:
		return (e.Queue[e.LastFloor][def.ButtonDown] == 1) || (e.Queue[e.LastFloor][def.ButtonInternal] == 1) || !requestsBelow(e) || e.LastFloor == 0
	case def.DirUp:
		return (e.Queue[e.LastFloor][def.ButtonUp] == 1) || (e.Queue[e.LastFloor][def.ButtonInternal] == 1) || !requestsAbove(e) || e.LastFloor == 3
	case def.DirStop:
	default:
		return true
	}
	return true
}

func Enqueue(e *def.Elevator, order def.Order) {
	e.Queue[order.Floor][order.Type] = 1
	backup.BackupInternalQueue(*e)
}

func UpdateGlobalQueue(globalQueue_chan chan [4][2]int, old_queue [4][2]int, newOrder def.Order) {
	if newOrder.Type == def.ButtonInternal {
		globalQueue_chan <- old_queue
	} else {
		old_queue[newOrder.Floor][int(newOrder.Type)] = 1
		globalQueue_chan <- old_queue
	}
}
