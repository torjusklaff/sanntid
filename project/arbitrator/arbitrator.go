package arbitrator

import (
	"math"
	def "../definitions"
)

var max_distance int = def.N_floors * def.N_buttons

func Find_lowest_cost(IP_adresses []string, costs []float32) to_elevator {
	lowest_cost := math.Inf(+1)
	to_elevator string := "lol"
	for i := 0; i<len(costs); i++{
		if costs[i] < lowest_cost {
			lowest_cost = costs[i]
			to_elevator = IP_adresses[i]
		}
	}
}

func Cost_function(elevator def.Elevator, order def.Order_button) cost{
	difference := order.Floor - elevator.Last_floor
	cost := math.Abs(difference) 
	+ movement_penalty(elevator.State, elevator.Current_direction, difference) 
	+ turn_penalty(elevator.State, elevator.Last_floor, elevator.Current_direction, order.Floor)
	+ order_direction_penalty(elevator.Current_direction, order.Floor, order.Type)
}


func movement_penalty(state def.Elev_states, direction def.Motor_direction, difference int) penalty{
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

func turn_penalty(state def.Elev_states, elevator_floor int, elevator_direction def.Motor_direction, order_floor int) penalty{
	if((state == idle)&&((elevator_floor == 1)||(elevator_floor == driver.N_floors)))||((state == moving)&&((first)||(second))) {
		penalty = 0
	} else if (elevator_direction==dir_up && order_floor<elevator_floor)|| (elevator_direction == dir_down && order_floor > elevator_floor){
		penalty = 0.75
	} else {
		penalty = 0
	}
}


func order_direction_penalty(elevator_direction def.Motor_direction, order_floor int, order_direction def.Motor_direction) penalty{
	if (order_floor == 1 || order_floor == driver.N_floors){
		penalty = 0
	} else if (elevator_direction != order_direction){
		penalty = driver.N_floors-2+0.25
	} else {
		penalty = 0
	}
}
