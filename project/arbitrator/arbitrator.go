package arbitrator

import (
	def "../definitions"
	//"fmt"
	"math"
	"strings"
)

var max_distance int = def.N_floors * def.N_buttons

func cost_function(e def.Elevator, new_order def.Order) int {
	cost_value := 0

	diff := new_order.Floor - e.Last_floor


	//Turn reward
	if (e.Current_direction == 1 && diff > 0) || (e.Current_direction == -1 && diff < 0) {
		cost_value -= 50
	//Turn penalty
	} else if (e.Current_direction == 1 && diff < 0) || (e.Current_direction == -1 && diff > 0) {
		cost_value += 225

	//Distance
	} else if (e.Current_direction == 0) {
		if (diff == -1) || (diff == 1) {
			cost_value += 25		
		} else if (diff == -2) || (diff == 2) {
			cost_value += 50
		} else if (diff == -3) || (diff == 3) {
			cost_value += 75
		}
	} 

	return cost_value
}

func arbitrator_optimal_next_order() {
	//enten lag eller ta inn en liste med alle bestillinger, send top til fsm_next_order
}

func find_lowest_cost(costs [def.N_elevators]def.Cost) def.Cost {
	for i := 0; i < len(costs)-1; i++ {
		if costs[i+1].Cost < costs[i].Cost {
			temp := costs[i]
			costs[i] = costs[i+1]
			costs[i+1] = temp
		}
		if costs[0] == costs[1] {
			if split_IP(costs[0].Id) < split_IP(costs[1].Id) {
				return costs[0]
			} else {
				return costs[1]
			}
		}
	}
	return costs[0]
}

// initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
func Arbitrator_init(
	e def.Elevator,
	localIP string,
	receive_new_order chan def.Order,
	assigned_new_order chan def.Order,
	receive_cost chan def.Cost,
	send_cost chan def.Cost,
	number_of_connected_elevators chan int,
	handled_order chan int) {

	var n_elevators int
	State_controller := make([n_elevators]def.Elev_states)

	for {
		select {
		case elevators := <-number_of_connected_elevators:
			n_elevators = elevators
		case state_update := <-receive_cost:
			State_controller[state_update.Id] = state_update
		case current_new_order := <-receive_new_order:
			current_cost := def.Cost{Cost: cost_function(e, current_new_order), Current_order: current_new_order, Id: localIP}
			order_selection(e, assigned_new_order, State_controller, n_elevators, current_new_order, localIP)
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assigned_new_order
func order_selection(e def.Elevator,
	assigned_new_order chan<- def.Order,
	State_controller []def.Elev_states,
	n_elevators int,
	current_order def.Order,
	localIP string) {

	const n = n_elevators
	var cost_list []def.Cost
	i := 0
	for index, state := range State_controller{
		cost_value := cost_function(e, current_order)
		current_order.Id = index
		cost := def.Cost{cost_value, current_order}
		cost_list[i] = cost_value
		i++
	}

	lowest_cost := math.Inf(+1)

	for index, cost := range cost_list {
		if cost_list[index] < lowest_cost {
			lowest_cost = cost_list[index]
			id := index
		}
	}

	if id == localIP{
		assigned_new_order <- lowest_cost.Current_order
	}
}

//hjelpefunksjon
func split_IP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

