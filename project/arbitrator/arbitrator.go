package arbitrator

import (
	"math"
	def "../definitions"
	"strings"
)

var max_distance int = def.N_floors * def.N_buttons


func cost_function(elevator def.Elevator, order def.Order_button) cost{
	difference := order.Floor - elevator.Last_floor
	cost := math.Abs(difference) 
	+ movement_penalty(elevator.State, elevator.Current_direction, difference) 
	+ turn_penalty(elevator.State, elevator.Last_floor, elevator.Current_direction, order.Floor)
	+ order_direction_penalty(elevator.Current_direction, order.Floor, order.Type)
}

func find_lowest_cost(costs [def.N_elevators]def.Cost) def.Cost {
	for i := 0; i<len(costs)-1; i++{
		if costs[i+1].Cost_value < costs[i].Cost_value {
			temp := costs[i]
			costs[i] = costs[i+1]
			costs [i+1] = temp
		}
		if costs[0] == costs[1] {
			if (split_IP(costs[0].Id) < split_IP(costs[1].Id)){
				return list[0]
			} else {
				return list[1]
			}
		}
	}
	return list[0]
}


func Arbitrator_init(
	e def.Elevator,
	localIP string, 
	new_order <-chan def.Order_button, 
	assigned_new_order <- chan def.Order_button,
	receive_cost <-chan def.Cost, 
	send_cost chan<- def.Cost, 
	number_of_connected_elevators <-chan int){

	var n_elevators int


	for {
		select {
			case elevators := <- number_of_connected_elevators:
				n_elevators = elevators
			case current_new_order := <- new_order:
				current_cost := def.Cost{Cost: cost_function(e, current_new_order), Current_order: current_new_order, Id: localIP}
				order_selection(assigned_new_order, receive_cost, n_elevators, current_cost, localIP)
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assigned_new_order
func order_selection(
	assigned_new_order chan<- def.Order_button, 
	receive_cost <- chan def.Cost, 
	n_elevators int, 
	current_cost def.Cost, 
	localIP string){

	var cost_list [def.N_elevators]def.Cost

	for i := 0; i < def.N_elevators; i++ {
		cost_list[i] := def.Cost{math.Inf(+1), current_cost.Current_order, current_cost.Id}
	}

	switch (n_elevators){
	case 1:
		cost_list[0] = cost
	case 2:
		new_cost := <-elev_receive_cost_value
		cost_list[0] = cost
		cost_list[1] = new_cost
	case 3:
		new_cost := <-elev_receive_cost_value
		new_cost2 := <-elev_receive_cost_value
		cost_list[0] = cost
		cost_list[1] = new_cost
		cost_list[2] = new_cost2
	}

	// regner ut laveste kost av de aktive heisene
	lowest_cost := find_lowest_cost(cost_list)

	// sender 
	if lowest_cost.Id == localIP {
		assigned_new_order <- cost.Current_order
	}
}






func split_IP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
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
