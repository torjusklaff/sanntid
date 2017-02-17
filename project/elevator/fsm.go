package fsm 

import (
	"../driver"
	"../elevator"
	"../queue"
)

elevator elevator.Elevator

// What is this for?
func set_all_lights(e Elevator){
	for floor := 0; floor < N_floors; floor++{
		for btn := 0; btn < N_buttons; btn++{
			driver.Set_button_lamp(btn, floor, 1)
		}
	}
}


func FSM_init_between_floors(){
	driver.Set_motor_direction(dir_down)
	elevator.current_direction = dir_down
	elevator.elevator_state = moving
}

/*
func fsm_button_pressed(btn_floor int, btn_type Button_type){
	switch(elevator.elevator_state){
	case door_open:
		if(elevator.last_floor == btn_floor){
			timer_start(elevator.door_open_duration)
		}
		else {
			elevator.queue[btn_floor][btn_type] = 1
		}
		break
	case moving:
		elevator.queue[btn][btn_type] = 1
	case idle:
	
	}
}
*/ //in arbitrator


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
			//set_all_lights(elevator)  			// why?
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
		elevator.current_direction = queue.choose_direction(elevator)

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




