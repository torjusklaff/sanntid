package arbitrator

import (
	def "../definitions"
	"fmt"
	"math"
	"strings"
)

var max_distance int = def.N_floors * def.N_buttons

func costFunction(elevator def.Elevator, order def.Order) float64 {
	difference := order.Floor - elevator.Last_floor
	cost := math.Abs(float64(difference)) + movementPenalty(elevator.Elevator_state, elevator.Current_direction, difference) + turnPenalty(elevator.Elevator_state, elevator.Last_floor, elevator.Current_direction, order.Floor) + orderDirectionPenalty(elevator.Current_direction, order.Floor, order.Type)
	return cost
}

func arbitratorOptimalNextOrder() {
	//enten lag eller ta inn en liste med alle bestillinger, send top til fsm_next_order
}

func findLowestCost(costs [def.N_elevators]def.Cost) def.Cost {
	for i := 0; i < len(costs)-1; i++ {
		if costs[i+1].Cost < costs[i].Cost {
			temp := costs[i]
			costs[i] = costs[i+1]
			costs[i+1] = temp
		}
		if costs[0] == costs[1] {
			if splitIP(costs[0].Id) < splitIP(costs[1].Id) {
				return costs[0]
			} else {
				return costs[1]
			}
		}
	}
	return costs[0]
}

// initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
func ArbitratorInit(
	e def.Elevator,
	localIP string,
	receive_new_order chan def.Order,
	assigned_new_order chan def.Order,
	received_states chan def.Elevator,
	send_states chan def.Elevator,
	number_of_connected_elevators chan int) {

	elev_states := make(map[string][]def.Elevator)
	n_elevators := 1

	for {
		select {
		case elevators := <-number_of_connected_elevators:
			n_elevators = elevators
			fmt.Printf("Number of elevators: %v \n", n_elevators)
		case current_new_order := <-receive_new_order:
			fmt.Printf("We receive a new order\n")
			send_states <- e
			if (current_new_order.Type == def.Buttoncall_internal) || (n_elevators == 1) {
				assigned_new_order <- current_new_order
			} else {
				//current_cost := def.Cost{Cost: costFunction(e, current_new_order), Current_order: current_new_order, Id: localIP}
				//orderSelection(assigned_new_order, receive_cost, n_elevators, current_cost, localIP)
				//if current_new_order.Floor == 3 || current_new_order.Type == def.Buttoncall_up {
				//	assigned_new_order <- current_new_order
			}
		case new_states := <-received_states:
			elev_states[new_states.Id] = append(elev_states[new_states.Id], new_states)
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assigned_new_order
func orderSelection(
	assigned_new_order chan<- def.Order,
	receive_cost <-chan def.Cost,
	n_elevators int,
	current_cost def.Cost,
	localIP string) {

	var cost_list [def.N_elevators]def.Cost

	for i := 0; i < def.N_elevators; i++ {
		cost_list[i] = def.Cost{math.Inf(+1), current_cost.Current_order, current_cost.Id}
	}

	/*switch n_elevators {
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
	}*/

	// regner ut laveste kost av de aktive heisene
	lowest_cost := findLowestCost(cost_list)

	// sender
	if lowest_cost.Id == localIP {
		assigned_new_order <- current_cost.Current_order
		fmt.Printf("We took the order!\n")
	} else {
		fmt.Printf("Someone else took the order\n")
	}
}

//hjelpefunksjon for å velge hvis cost er lik
func splitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func movementPenalty(state def.Elev_states, direction def.Motor_direction, difference int) float64 {
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

func turnPenalty(state def.Elev_states, elevator_floor int, elevator_direction def.Motor_direction, order_floor int) float64 {
	if ((state == def.Idle) && ((elevator_floor == 1) || (elevator_floor == def.N_floors))) || (state == def.Moving) {
		return 0
	} else if (elevator_direction == def.Dir_up && order_floor < elevator_floor) || (elevator_direction == def.Dir_down && order_floor > elevator_floor) {
		return 0.75
	} else {
		return 0
	}
}

func orderDirectionPenalty(elevator_direction def.Motor_direction, order_floor int, order_direction def.Button_type) float64 {
	if order_floor == 1 || order_floor == def.N_floors {
		return 0
	} else if int(elevator_direction) != int(order_direction) {
		return def.N_floors - 2 + 0.25
	} else {
		return 0
	}
}
