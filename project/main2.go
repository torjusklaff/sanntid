package main

// Test-main for driver-files
import (
	//arb "./arbitrator"
	//"./driver"
	"./backup"
	def "./definitions"
	//"./fsm"
	//net "./network"
	q "./queue"
	"fmt"
	"time"
	"os"
	"os/signal"
	"log"
)

func main() {

	door_timer := time.NewTimer(3 * time.Second)
	door_timer.Stop()

	var elevator def.Elevator
	elevator.Last_floor = 1
	elevator.Current_direction = def.Dir_stop
	elevator.Queue = [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	elevator.Elevator_state = def.Idle

	var previous_order def.Order
	previous_order.Type = def.Buttoncall_internal
	previous_order.Floor = elevator.Last_floor
	previous_order.Id = "-"
	previous_order.Internal = true


	// 		BACKUP KAN NÅ LAGRE TING I FIL, SAMT AT KØ-MODULEN KAN DECODE STRINGS TIL KØ-ARRAYS
	queue_string := q.Queue_to_string(elevator)
	backup.To_backup(queue_string)

	string_size := len(queue_string)
	last_line := backup.Read_last_line(int64(string_size))
	fmt.Print(last_line)

	go func(){
		var c = make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Print("User terminated program.\n")
		var err = os.Remove("log.txt")
		if err != nil {
	        log.Fatalf("Error deleting file: %v", err)
	    }
		log.Fatal("User terminated program.\n")
	}()

	for{}
	//queue := q.Queue_from_string(last_line+"\n")
	


	// 	CHANNELS 
	/*
	n_elevators := make(chan int)

	receive_cost := make(chan def.Cost)
	receive_new_order := make(chan def.Order)
	receive_remove_order := make(chan def.Order)

	send_cost := make(chan def.Cost)
	send_new_order := make(chan def.Order)
	send_remove_order := make(chan def.Order)
	assigned_new_order := make(chan def.Order)

	button_pressed := make(chan def.Order)
	on_floor := make(chan int)

	id := net.Get_id()
	go net.Network_init(id, n_elevators, receive_cost, receive_new_order, receive_remove_order, send_cost, send_new_order, send_remove_order)
	go arb.Arbitrator_init(elevator, id, receive_new_order, assigned_new_order, receive_cost, send_cost, n_elevators) // button_pressed må endres til receive_new_order

	go driver.Check_all_buttons(send_new_order)
	go driver.Elevator_on_floor(on_floor, elevator)


	for {
		select {
		case floor := <-on_floor:
			fsm.FSM_floor_arrival(floor, &elevator, door_timer)
		case <-door_timer.C:
			fmt.Printf("Timer stopped\n")
			fsm.FSM_on_door_timeout(&elevator)
		case new_order := <-assigned_new_order:
			fmt.Print("Assigned new order\n")
			queue.Enqueue(&elevator, new_order)
			fsm.FSM_next_order(&elevator, new_order, door_timer)
		default:
			break
		}
	}
	*/
}


