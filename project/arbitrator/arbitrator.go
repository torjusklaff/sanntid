package arbitrator

import (
	def "../definitions"
	"fmt"
	"math"
	"strings"
)


type reply struct{
	cost int
	lift string
}
type order struct {
	floor  int
	button int
	timer  *time.Timer
}


func Arbitrator.run(cost_reply chan def.Cost, n_elevators *int){
	unassigned := make(map[order][]reply)
	var timeout = make(chan *def.Order)

	for{
		select{
		case message := <- cost_reply:
			new_order := order{Floor: message.Floor, Type: message.Type}
			new_reply := reply{cost: message.Cost, lift: message.Id}

			for old_order := range unassigned{
				if equal(old_order, new_order){
					new_order = old_order
				}
			}

			if reply_list, exist := unassigned[new_order]; exist{
				found := false
				for _, reply := range reply_list{
					if reply == new_reply{
						found = true
					}
				}
				if !found {
					unassigned[new_order] = append(unassigned[new_order], new_reply)
					new_order.timer.Reset(def.Order_time_out)
				}
			} else {
				new_order.timer = time.NewTimer(def.Order_time_out)
				unassigned[new_order] = []reply{new_reply}
				go cost_timer(&new_order, timeout)
			}
			choose_best_lift(unassigned, n_elevators, false)
		case <-timeout:
			log.Println("Not all costs received in time")
			choose_best_lift(unassigned, n_elevators, true)
		}
	}
}

func choose_best_lift(unassigned map[order][]reply, n_elevators *int, order_timed_out bool){
	const max_int = int(^uint(0) >> 1)
	// Loop through all lists.
	for order, reply_list := range unassigned {
		// Check if the list is complete or the timer has timed out.
		if len(reply_list) == *n_elevators || order_timed_out {
			lowestCost := maxInt
			var best_lift string

			// Loop through costs in each complete list.
			for _, reply := range reply_list {
				if reply.cost < lowest_cost {
					lowest_cost = reply.cost
					best_lift = reply.lift
				} else if reply.cost == lowest_cost {
					// Prioritise on lowest IP value if cost is the same.
					if reply.lift < best_lift {
						lowest_cost = reply.cost
						best_lift = reply.lift
					}
				}
			}
			queue.AddRemoteOrder(order.floor, order.button, bestLift)
			order.timer.Stop()
			delete(unassigned, order)
		}
	}
}

func cost_timer(newOrder *order, timeout chan<- *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
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
			current_cost := def.Cost{Cost: cost_function(e, current_new_order), Current_order: current_new_order, Id: localIP}
			order_selection(assigned_new_order, receive_cost, n_elevators, current_cost, localIP)
		}
	}
}
