package queue

import (
	def "../definitions"
	"fmt"
	//"log"
	"../backup"
)

func requests_above(e def.Elevator) bool {
	for f := e.Last_floor + 1; f < def.N_floors; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {
			if e.Queue[f][def.Button_type(btn)] == 1 {
				return true
			}
		}
	}
	return false
}

func requests_below(e def.Elevator) bool {
	for f := 0; f < e.Last_floor; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {
			if e.Queue[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func Choose_direction(e def.Elevator) def.Motor_direction {
	switch e.Current_direction {
	case def.Dir_up:
		if requests_above(e) {
			return def.Dir_up
		} else if requests_below(e) {
			return def.Dir_down
		} else {
			return def.Dir_stop
		}
	case def.Dir_down:
		if requests_below(e) {
			return def.Dir_down
		} else if requests_above(e) {
			return def.Dir_up
		} else {
			return def.Dir_stop
		}
	case def.Dir_stop:
		if requests_below(e) {
			return def.Dir_down
		} else if requests_above(e) {
			return def.Dir_up
		} else {
			return def.Dir_stop
		}
	default:
		return def.Dir_stop
	}
	return def.Dir_stop
}

func Clear_at_floor(e *def.Elevator, floor int) {
	for btn := 0; btn < def.N_buttons; btn++ {
		if e.Queue[floor][btn] == 1{
			e.Queue[floor][btn] = 0
			backup.Backup_internal_queue(*e)
		}
	}
}

func Clear_global_queue(send_global_queue chan [4][2]int, old_queue [4][2]int, floor int){
	for btn := 0; btn < 2; btn++ {
		if old_queue[floor][btn] == 1{
			old_queue[floor][btn] = 0
		}
	}
	send_global_queue <- old_queue
}


func Print_queue(e def.Elevator) {
	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {
			fmt.Printf("%v ", e.Queue[f][btn])
		}
		fmt.Printf("\n")
	}
}



func Should_stop(e def.Elevator) bool {
	switch e.Current_direction {
	case def.Dir_down:
		return (e.Queue[e.Last_floor][def.Buttoncall_down] == 1) || (e.Queue[e.Last_floor][def.Buttoncall_internal] == 1) || !requests_below(e)
	case def.Dir_up:
		return (e.Queue[e.Last_floor][def.Buttoncall_up] == 1) || (e.Queue[e.Last_floor][def.Buttoncall_internal] == 1) || !requests_above(e)
	case def.Dir_stop:
	default:
		return true
	}
	return true
}


func Enqueue(e *def.Elevator, order def.Order) {
	e.Queue[order.Floor][order.Type] = 1
	backup.Backup_internal_queue(*e)
}

func Update_global_queue(global_queue_chan chan [4][2]int, old_queue [4][2]int, new_order def.Order){
	if new_order.Type == def.Buttoncall_internal{
		global_queue_chan <- old_queue
	} else {
		old_queue[new_order.Floor][int(new_order.Type)] = 1
		global_queue_chan <- old_queue
	}
}

