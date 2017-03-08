package fsm 

import (
	"../driver"
	"../elevator"
	"../queue"
)

elevator elevator.Elevator




func FSM_floor_arrival(new_floor int){
	Elevator.last_floor = new_floor
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


func FSM_on_door_timeout(){
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