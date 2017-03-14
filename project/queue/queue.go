package queue

import (
	def "../definitions"
	"../backup"
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

func DeleteInternalQueuesAtFloor(e *def.Elevator, floor int) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if e.Queue[floor][btn] == 1 {
			e.Queue[floor][btn] = 0
			backup.BackupInternalQueue(*e)
		}
	}
}

func DeleteGlobalQueuesAtFloor(SendGlobalQueue chan [4][2]int, oldQueue [4][2]int, floor int) {
	for btn := 0; btn < 2; btn++ {
		if oldQueue[floor][btn] == 1 {
			oldQueue[floor][btn] = 0
		}
	}
	SendGlobalQueue <- oldQueue
}


func ShouldStop(e def.Elevator) bool {
	switch e.CurrentDirection {
	case def.DirDown:
		return (e.Queue[e.LastFloor][def.ButtoncallDown] == 1) || (e.Queue[e.LastFloor][def.ButtoncallInternal] == 1) || !requestsBelow(e) || e.LastFloor == 0
	case def.DirUp:
		return (e.Queue[e.LastFloor][def.ButtoncallUp] == 1) || (e.Queue[e.LastFloor][def.ButtoncallInternal] == 1) || !requestsAbove(e) || e.LastFloor == 3
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

func AddOrderToGlobalQueue(globalQueueChan chan [4][2]int, oldQueue [4][2]int, newOrder def.Order) {
	if newOrder.Type == def.ButtoncallInternal {
		globalQueueChan <- oldQueue
	} else {
		oldQueue[newOrder.Floor][int(newOrder.Type)] = 1
		globalQueueChan <- oldQueue
	}
}

