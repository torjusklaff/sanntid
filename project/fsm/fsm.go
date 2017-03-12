package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	"fmt"
	"time"
)

func FSM_floor_arrival(new_floor int, elevator *def.Elevator) {
	if new_floor == -1 {
		fmt.Print("Run FSM_floor_arrival while not on floor\n")
	} else if (new_floor == elevator.Last_floor){
	} else {
		fmt.Print("FSM_floor_arrival\n")
		driver.Set_floor_indicator(new_floor)
		elevator.Last_floor = new_floor
		elevator.Motor_stop_timer.Stop()
		switch elevator.Elevator_state {
		case def.Moving:
			if queue.Should_stop(*elevator) {
				driver.Set_motor_direction(def.Dir_stop)
				queue.Clear_at_floor(elevator, new_floor)
				driver.Clear_lights_at_floor(elevator.Last_floor)

				driver.Set_door_open_lamp(1)
				elevator.Door_timer.Reset(3 * time.Second)
				fmt.Printf("Timer started\n")
				elevator.Elevator_state = def.Stop_on_floor
			}
			break
		case def.Idle:
		default:
			break
		}
	}
}



func FSM_next_order(elevator *def.Elevator, next_order def.Order) { //arbitrator decides where we should go next
	fmt.Print("FSM_next_order\n")
	driver.Set_button_lamp(next_order, 1)
	
	switch elevator.Elevator_state {
	case def.Idle:
		queue.Enqueue(elevator, next_order)
		if next_order.Floor == elevator.Last_floor {
			queue.Clear_at_floor(elevator, elevator.Last_floor)
			driver.Clear_lights_at_floor(elevator.Last_floor)
			elevator.Door_timer.Reset(3 * time.Second)
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
			elevator.Motor_stop_timer.Reset(4*time.Second)
			fmt.Print("FSM_next_order: Reset motor_timer\n")
		}
	case def.Moving:
	case def.Stop_on_floor:
		queue.Clear_at_floor(elevator, elevator.Last_floor)
		driver.Clear_lights_at_floor(elevator.Last_floor)
	case def.Motor_stop:
		if next_order.Type == def.Buttoncall_internal{
			queue.Enqueue(elevator, next_order)
		}
			
	default:
		break
	}
	queue.Print_queue(*elevator)
}

func FSM_on_door_timeout(elevator *def.Elevator) {
	fmt.Print("FSM_on_door_timeout\n")
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
			elevator.Motor_stop_timer.Reset(4*time.Second)
			fmt.Print("FSM_where_to_next: Reset motor_timer\n")
		}
		break
	case def.Idle:
		elevator.Current_direction = queue.Choose_direction(*elevator)
		driver.Set_motor_direction(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(4*time.Second)
			fmt.Print("FSM_on_door_timeout: Reset motor_timer\n")
		}
		break
	default:
		break
	}
}

func FSM_where_to_next(elevator def.Elevator){
	switch elevator.Elevator_state {
	case def.Stop_on_floor:
		elevator.Current_direction = queue.Choose_direction(elevator)
		driver.Set_motor_direction(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(4*time.Second)
			fmt.Print("FSM_where_to_next: Reset motor_timer\n")
		}
		break
	case def.Idle:
		elevator.Current_direction = queue.Choose_direction(elevator)
		driver.Set_motor_direction(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(4*time.Second)
			fmt.Print("FSM_on_door_timeout: Reset motor_timer\n")
		}
		break
	default:
		break
	}
}


func FSM_motor_stop(elevator *def.Elevator){
	fmt.Print("FSM_motor_stop\n")
	elevator.Current_direction = def.Dir_stop
	driver.Set_motor_direction(def.Dir_stop)

	driver.Elev_init_from_backup()

	/*dead := true
	for dead{
		driver.Set_motor_direction(def.Dir_down)
		if driver.Get_floor_sensor_signal() != -1 {
			fmt.Print(int(driver.Get_floor_sensor_signal()))
			driver.Set_motor_direction(def.Dir_stop)
			dead = false
		}

	}*/
}


