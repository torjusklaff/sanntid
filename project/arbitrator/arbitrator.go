package arbitrator

import (
	def "../definitions"
	"fmt"
	"math"
	"strings"
)

var max_distance int = def.N_floors * def.N_buttons

// initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
func ArbitratorInit(
	e def.Elevator,
	receive_new_order chan def.Order,
	assigned_new_order chan def.Order,
	//received_states chan def.Elevator,
	number_of_connected_elevators chan int,
	elev_states map[string]def.Elevator_msg) {

	n_elevators := 1
	costs := make(map[string]def.Cost)
	for {
		select {
		case elevators := <-number_of_connected_elevators:
			n_elevators = elevators
			fmt.Printf("Number of elevators: %v \n", n_elevators)
		case current_new_order := <-receive_new_order:
			if n_elevators == 1 {
				fmt.Printf("We are alone, we get the order!\n")
				assigned_new_order <- current_new_order
			} else {
				/*
					new_states := <-received_states
					elev_states[e.Id]= e
					elev_states[new_states.Id] = new_states
				*/ //FORSLAG: fjerne disse linjene da de kan skape en sleep
				for elevator_id, _ := range elev_states {
					costs[elevator_id] = def.Cost{Cost: costFunction(elev_states[elevator_id], current_new_order), Current_order: current_new_order, Id: elevator_id}
				}
				orderSelection(assigned_new_order, costs, n_elevators, e.Id)

			}
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assigned_new_order
func orderSelection(
	assigned_new_order chan<- def.Order,
	cost_list map[string]def.Cost,
	n_elevators int,
	localIP string) {
	lowest_cost := findLowestCost(cost_list)
	fmt.Printf("Lowest cost calculated\n")
	// sender
	if lowest_cost.Id == localIP {
		fmt.Printf("We took the order!\n")
		assigned_new_order <- lowest_cost.Current_order

	} else {
		fmt.Printf("Someone else took the order\n")
	}
}

func costFunction(e def.Elevator_msg, order def.Order) float64 {
	diff := order.Floor - e.Last_floor
	cost := math.Abs(float64(diff)) + movementPenalty(e.Elevator_state, e.Current_direction, diff) + turnPenalty(e.Elevator_state, e.Last_floor, e.Current_direction, order.Floor) + orderDirectionPenalty(e.Current_direction, order.Floor, order.Type)
	return cost
}

func findLowestCost(costs map[string]def.Cost) def.Cost { //Problemet er inni her!!!!!!!!!!
	dummy_order := def.Order{Type: 0, Floor: 0, Internal: false, Id: " "}
	lowest_cost := def.Cost{Cost: math.Inf(+1), Current_order: dummy_order, Id: "i am stupid"}

	for Id, cost := range costs {
		if cost.Cost < lowest_cost.Cost {
			lowest_cost = cost
		}
		if cost.Cost == lowest_cost.Cost {
			if splitIP(Id) < splitIP(lowest_cost.Id) {
				lowest_cost = cost
			}
		}
	} //Skjer i for løkka at det blir index out of range. FOr sliten til å fikse nå, se gjerne på det selv hvis du får to heiser
	return lowest_cost
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
