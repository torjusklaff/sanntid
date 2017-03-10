package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	//arb "../arbitrator"
	"fmt"
	"time"
)

func FSM_floor_arrival(new_floor int, elevator *def.Elevator, timer *time.Timer) {
	if new_floor == -1 {
	} else {
		driver.Set_floor_indicator(new_floor)
		elevator.Last_floor = new_floor

		switch elevator.Elevator_state {
		case def.Moving:
			if queue.Should_stop(*elevator) {
				driver.Set_motor_direction(def.Dir_stop)
				queue.Clear_at_floor(elevator, new_floor)
				driver.Clear_lights_at_floor(elevator.Last_floor)

				driver.Set_door_open_lamp(1)
				timer.Reset(3 * time.Second)
				fmt.Printf("Timer started\n")
				elevator.Elevator_state = def.Stop_on_floor
			}
			break
		default:
			break
		}
	}
}

func FSM_next_order(elevator *def.Elevator, next_order def.Order, timer *time.Timer) { //arbitrator decides where we should go next
	driver.Set_button_lamp(next_order, 1)
	switch elevator.Elevator_state {
	case def.Idle:

		if next_order.Floor == elevator.Last_floor {
			queue.Clear_at_floor(elevator, elevator.Last_floor)
			driver.Clear_lights_at_floor(elevator.Last_floor)
			timer.Reset(3 * time.Second)
			driver.Set_door_open_lamp(1)
			elevator.Elevator_state = def.Stop_on_floor
		} else {

			if next_order.Floor > elevator.Last_floor {
				elevator.Current_direction = def.Dir_up
				driver.Set_motor_direction(elevator.Current_direction)
			} else {
				elevator.Current_direction = def.Dir_down
				driver.Set_motor_direction(elevator.Current_direction)
			}

		}
		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
		}
	case def.Moving:
	case def.Stop_on_floor:
		queue.Clear_at_floor(elevator, elevator.Last_floor)
		driver.Clear_lights_at_floor(elevator.Last_floor)
	default:
		break
	}
}

func FSM_on_door_timeout(elevator *def.Elevator) {
	queue.Print_queue(*elevator)
	driver.Set_door_open_lamp(0)
	switch elevator.Elevator_state {
	case def.Stop_on_floor:
		elevator.Current_direction = queue.Choose_direction(*elevator)
		driver.Set_motor_direction(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
		}

		break
	default:
		break
	}
}
