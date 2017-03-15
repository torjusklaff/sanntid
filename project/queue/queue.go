package queue

import (
	"../backup"
	def "../definitions"
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

func UpdateGlobalQueue(globalQueue *[4][2]int, order def.Order) {
	globalQueue[order.Floor][int(order.Type)] = 1
}
func DeleteGlobalOrdersAtFloor(globalQueue *[4][2]int, floor int) {
	for i := 0; i < def.NumButtons-1; i++ {
		globalQueue[floor][i] = 0
	}
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

func DeleteInternalQueueAtFloor(e *def.Elevator, floor int) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if e.Queue[floor][btn] == 1 {
			e.Queue[floor][btn] = 0
			backup.BackupInternalQueue(*e)
		}
	}
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
