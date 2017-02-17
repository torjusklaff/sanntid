package main

import (
	"/driver"
	"/elevator"
)

func requests_above(e Elevator) int {
	for int f := e.last_floor+1; f < N_floors; f++ {
		for btn := 0; btn < N_buttons; btn++ {
			if e.queue[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e Elevator) int {
	for int f := 0; f < e.last_floor; f++ {
		for btn := 0; btn < N_buttons; btn++ {
			if e.queue[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func Choose_direction(e Elevator) motor_direction {
	switch(e.current_direction){
	case dir_up:
		if requests_above(e){
			return dir_up
		}
		if else requests_below(e){
			return dir_down
		}
		else {
			return dir_stop
		}
   case dir_down:
   case dir_stop:
   		if requests_below(e){
			return dir_down
		}
		if else requests_above(e){
			return dir_up
		}
		else {
			return dir_stop
		}
   	default:
   		return dir_stop
	}
}

func Clear_at_current_floor(e Elevator){
	for btn Button_type := 0; btn < N_buttons; btn++{
		e.queue[e.last_floor][btn] = 0;
	}
}

func Should_stop(e Elevator) int{
	switch(e.current_direction){
	case dir_down:
		return e.queue[e.last_floor][Buttoncall_down] || e.queue[e.floor][Buttoncall_internal] || !requests_below(e)
	case dir_up:
		return e.queue[e.last_floor][Buttoncall_up] || e.queue[e.floor][Buttoncall_internal] || !requests_above(e)
	case dir_stop:
	default:
		return 1
	}
}

func Enqueue(e Elevator, order Order_button){
	e.queue[order.Floor][order.Type] = 1;
}




