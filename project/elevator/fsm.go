package fsm 

import (
	"/driver"
	def "/definitions"
	"/queue"
	arb "/arbitrator"
)

func FSM_button_pressed(button def.Order_button, elevator def.Elevator) arbitrator_cost int{
	driver.Set_button_lamp(button.Type, button.Floor, 1)
	arbitrator_cost = arb.Cost_function(elevator, button)
}


func FSM_floor_arrival(new_floor int, elevator def.Elevator){
	elevator.last_floor = new_floor
	driver.Set_floor_indicator(new_floor)

	switch(elevator.elevator_state){
	case moving:
		if(queue.Should_stop(elevator)){
			driver.Set_motor_direction(dir_stop)
			driver.Set_door_open_lamp(1)
			queue.Clear_at_current_floor(elevator)
			timer_start(elevator.door_open_duration)
			elevator.elevator_state = door_open
		}
		break
	default:
		break
	}
}


func FSM_on_door_timeout(elevator def.Elevator){
	switch(elevator.elevator_state){
	case door_open:
		elevator.current_direction = queue.Choose_direction(elevator)

		driver.Set_door_open_lamp(0)
		driver.Set_motor_direction(elevator.current_direction)

		if(elevator.current_direction == dir_stop){
			elevator.elevator_state = idle
		}
		else {
			elevator.elevator_state = moving
		}

		break
	default:
		break
	}
}