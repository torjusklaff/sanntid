package arbitrator

import (
	def "../definitions"
	//"fmt"
	"math"
	"strings"
)

var max_distance int = def.N_floors * def.N_buttons

func cost_function(elevator def.Elevator, order def.Order) float64 {
	difference := order.Floor - elevator.Last_floor
	cost := math.Abs(float64(difference)) + movement_penalty(elevator.Elevator_state, elevator.Current_direction, difference) + turn_penalty(elevator.Elevator_state, elevator.Last_floor, elevator.Current_direction, order.Floor) + order_direction_penalty(elevator.Current_direction, order.Floor, order.Type)
	return cost
}

func Arbitrator_optimal_next_order() def.Order

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
	number_of_connected_elevators chan int) {

	var n_elevators int

	for {
		select {
		case elevators := <-number_of_connected_elevators:
			n_elevators = elevators
		case current_new_order := <-receive_new_order:
			fmt.Print("Registered new order in arbitrator\n")
			current_cost := def.Cost{Cost: cost_function(e, current_new_order), Current_order: current_new_order, Id: localIP}
			order_selection(assigned_new_order, receive_cost, n_elevators, current_cost, localIP)
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assigned_new_order
func order_selection(
	assigned_new_order chan<- def.Order,
	receive_cost <-chan def.Cost,
	n_elevators int,
	current_cost def.Cost,
	localIP string) {

	var cost_list [def.N_elevators]def.Cost

	for i := 0; i < def.N_elevators; i++ {
		cost_list[i] = def.Cost{math.Inf(+1), current_cost.Current_order, current_cost.Id}
	}

	switch n_elevators {
	case 1:
		cost_list[0] = current_cost
	case 2:
		new_cost := <-receive_cost
		cost_list[0] = current_cost
		cost_list[1] = new_cost
	case 3:
		new_cost := <-receive_cost
		new_cost2 := <-receive_cost
		cost_list[0] = current_cost
		cost_list[1] = new_cost
		cost_list[2] = new_cost2
	}

	// regner ut laveste kost av de aktive heisene
	lowest_cost := find_lowest_cost(cost_list)

	// sender
	if lowest_cost.Id == localIP {
		assigned_new_order <- current_cost.Current_order
	} else {
		fmt.Printf("Someone else took the order\n")
	}
}

//hjelpefunksjon
func split_IP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func movement_penalty(state def.Elev_states, direction def.Motor_direction, difference int) float64 {
	switch state {
	case def.Idle:
		return 0
	default:
		switch direction {
		case def.Dir_up:
			if difference > 0 {
				return -0.5
			} else if direction < 0 {
				return 1.5
			}
		case def.Dir_down:
			if difference > 0 {
				return 1.5
			} else if difference < 0 {
				return -0.5
			}
		}
	}
	return 0
}

func turn_penalty(state def.Elev_states, elevator_floor int, elevator_direction def.Motor_direction, order_floor int) float64 {
	if ((state == def.Idle) && ((elevator_floor == 1) || (elevator_floor == def.N_floors))) || (state == def.Moving) {
		return 0
	} else if (elevator_direction == def.Dir_up && order_floor < elevator_floor) || (elevator_direction == def.Dir_down && order_floor > elevator_floor) {
		return 0.75
	} else {
		return 0
	}
}

func order_direction_penalty(elevator_direction def.Motor_direction, order_floor int, order_direction def.Button_type) float64 {
	if order_floor == 1 || order_floor == def.N_floors {
		return 0
	} else if int(elevator_direction) != int(order_direction) {
		return def.N_floors - 2 + 0.25
	} else {
		return 0
	}
}
