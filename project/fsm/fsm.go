package fsm 

import (
	"../driver"
	def "../definitions"
	"../queue"
	//arb "../arbitrator"
)

func FSM_button_pressed(button def.Order_button, elevator *def.Elevator) /*arbitrator_cost int*/{
	driver.Set_button_lamp(button, 1)
	//arbitrator_cost = arb.Cost_function(elevator, button)

	queue.Enqueue(elevator, button)

	switch(elevator.Elevator_state){
	case def.Idle:
		elevator.Current_direction = queue.Choose_direction(*elevator)
		driver.Set_motor_direction(elevator.Current_direction)

		if(elevator.Current_direction == def.Dir_stop){
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
		}
	case def.Moving:
	default:
		break
	}

}

func FSM_floor_arrival(new_floor int, elevator *def.Elevator){
	elevator.Last_floor = new_floor
	driver.Set_floor_indicator(new_floor)

	switch(elevator.Elevator_state){
	case def.Moving:
		if(queue.Should_stop(*elevator)){
			driver.Set_motor_direction(def.Dir_stop)
			queue.Clear_at_floor(elevator, new_floor)

			//skrur av lys for den bestillingen som slettes
			for btn := 0; btn < def.N_buttons; btn++{
				var button def.Order_button
				button.Type = def.Button_type(btn)
				button.Floor = elevator.Last_floor
				driver.Set_button_lamp(button, 0)
			}

			driver.Door_open()
			elevator.Elevator_state = def.Door_open
			FSM_on_door_timeout(elevator)
		}
		break
	default:
		break
	}
}


func FSM_on_door_timeout(elevator *def.Elevator){
	switch(elevator.Elevator_state){
	case def.Door_open:
		elevator.Current_direction = queue.Choose_direction(*elevator)

		driver.Set_motor_direction(elevator.Current_direction)

		if(elevator.Current_direction == def.Dir_stop){
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
		}

		break
	default:
		break
	}
}