package queue

import (
	def "../definitions"
	//"fmt"
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

/*func Queue_not_empty(queue_not_empty chan def.Order_button, e def.Elevator) {
	for {

		for f := 0; f < def.N_floors; f++ {
			for btn := 0; btn < def.N_buttons; btn++ {
				fmt.Printf("%v ", e.Queue[f][btn])
				if e.Queue[f][btn] == 1 {
					var order def.Order_button
					order.Type = def.Button_type(btn)
					order.Floor = f
					fmt.Printf("We come this far\n")
					queue_not_empty <- order
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n\n")
	}
}*/

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
		e.Queue[floor][btn] = 0
		var button def.Order_button
		button.Type = def.Button_type(btn)
		button.Floor = floor
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

func Enqueue(e *def.Elevator, order def.Order_button) {
	e.Queue[order.Floor][order.Type] = 1
}
