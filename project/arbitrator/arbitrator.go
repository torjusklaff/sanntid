package arbitrator

import (
	"math"
	"/driver"	// path to Order_button
	"/elevator" // path to Elevator-struct
)

max_distance = driver.N_floors * driver.N_buttons

func Find_lowest_cost(IP_adresses string[], costs int[]) IP {
	var lowest_cost = inf
	var to_elevator IP;
	for (i=0; i<len(costs); i++){
		if costs[i] < lowest_cost {
			lowest_cost = costs[i]
			to_elevator = IP_adresses[i]
		}
	}
}

func Cost_function(elevator Elevator, order Order_button) cost int{
	var difference = order.Floor - elevator.last_floor
	var cost = math.Abs(difference) 
	+ movement_penalty(elevator.state, elevator.current_direction, difference) 
	+ turn_penalty(elevator.state, elevator.last_floor, elevator.current_direction, order.Floor)
	+ order_direction_penalty(elevator.current_direction, order.Floor, order.Type)
}


func movement_penalty(state, direction, difference) penalty int{
	switch(state){
	case idle:
		penalty = 0
	default:
		switch(direction){
		case dir_up:
			if (difference > 0){
				penalty = -0.5
			} else if (direction < 0){
				penalty = 1.5
			}
		case dir_down:
			if (difference > 0){
				penalty = 1.5
			} else if (difference < 0){
				penalty = -0.5
			}
		}
	}
}

func turn_penalty(state, elevator_floor, elevator_direction, order_floor) penalty int{
	if((state == idle)&&((elevator_floor == 1)||(elevator_floor == driver.N_floors)))
	||((state == moving)&&((first)||(second)))
	{
		penalty = 0
	} else if (elevator_direction==dir_up && order_floor<elevator_floor)
	|| (elevator_direction == dir_down && order_floor > elevator_floor){
		penalty = 0.75
	} else {
		penalty = 0
	}
}


func order_direction_penalty(elevator_direction, order_floor, order_direction) penalty int{
	if (order_floor == 1 || order_floor == driver.N_floors){
		penalty = 0
	} else if (elevator_direction != order_direction){
		penalty = driver.N_floors-2+0.25
	} else {
		penalty = 0
	}
}
