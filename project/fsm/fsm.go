package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	//arb "../arbitrator"
	"time"
)

func FSM_button_pressed(button def.Order_button, elevator *def.Elevator) /*arbitrator_cost int*/ {
	//arbitrator_cost = arb.Cost_function(elevator, button)
	switch elevator.Elevator_state {
	case def.Idle:
		if button.Floor == elevator.Last_floor {
			queue.Clear_at_floor(elevator, elevator.Last_floor)
		} else {
			elevator.Current_direction = queue.Choose_direction(*elevator)
			driver.Set_motor_direction(elevator.Current_direction)
		}
		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
		}
	default:
		break
	}
}

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

				//skrur av lys for den bestillingen som slettes
				for btn := 0; btn < def.N_buttons; btn++ {
					var button def.Order_button
					button.Type = def.Button_type(btn)
					button.Floor = elevator.Last_floor
					driver.Set_button_lamp(button, 0)
				}
				driver.Set_door_open_lamp(1)
				timer.Reset(3 * time.Second)
				elevator.Elevator_state = def.Door_open
			}
			break
		default:
			break
		}
	}
}

func FSM_next_order(elevator *def.Elevator, next_order def.Order_button) { //arbitrator decides where we should go next
	switch elevator.Elevator_state {
	case def.Idle:
		if next_order.Floor == elevator.Last_floor {
			queue.Clear_at_floor(elevator, elevator.Last_floor)
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
	case def.Door_open:

	default:
		break
	}
}

func FSM_on_door_timeout(elevator *def.Elevator) {
	driver.Set_door_open_lamp(0)
	switch elevator.Elevator_state {
	case def.Door_open:
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

// Prøvde å lage en ny funksjon som kun kunne sjekke knapper hele tiden
func Button_listener(button_press chan<- def.Order_button) {
	possible_buttons := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	for {
		for floor := 0; floor < def.N_floors; floor++ {
			for btn := def.Buttoncall_down; int(btn) < def.N_buttons; btn++ {
				if floor == 0 && btn == def.Buttoncall_down {
					continue
				}
				if floor == def.N_floors-1 && btn == def.Buttoncall_up {
					continue
				}

				var button def.Order_button
				button.Type = btn
				button.Floor = floor
				button_signal := driver.Get_button_signal(button)

				if button_signal == 1 && (possible_buttons[floor][btn] == 0) {
					button_press <- def.Order_button{Type: btn, Floor: floor}
					possible_buttons[floor][btn] = driver.Get_button_signal(button)
					driver.Set_button_lamp(button, 1)
				}
			}
		}
	}
}

// Prøvde å hele tiden sjekke etter floor-signaler
func Floor_listener(floor_pass chan<- int) {
	last_floor := -1
	var floor_signal int
	for {
		floor_signal = driver.Get_floor_sensor_signal()
		if (floor_signal != -1) && (last_floor != floor_signal) {
			floor_pass <- floor_signal
			driver.Set_floor_indicator(floor_signal)
		}
	}
}
