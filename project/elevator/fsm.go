package fsm 

import (
	"/driver"
	def "/definitions"
	"/queue"
	arb "/arbitrator"
)

func FSMButtonPressed(button def.Order_button, elevator def.Elevator) arbitrator_cost int{
	driver.Set_button_lamp(button.Type, button.Floor, 1)
	arbitrator_cost = arb.Cost_function(elevator, button)
}


func FSMFloor_arrival(newFloor int, elevator def.Elevator){
	elevator.LastFloor = newFloor
	driver.SetFloor_indicator(newFloor)

	switch(elevator.ElevatorState){
	case moving:
		if(queue.Should_stop(elevator)){
			driver.Set_MotorDirection(DirStop)
			driver.Set_doorOpen_lamp(1)
			queue.Clear_at_currentFloor(elevator)
			timer_start(elevator.doorOpen_duration)
			elevator.ElevatorState = doorOpen
		}
		break
	default:
		break
	}
}


func FSMOn_door_timeout(elevator def.Elevator){
	switch(elevator.ElevatorState){
	case doorOpen:
		elevator.CurrentDirection = queue.Choose_direction(elevator)

		driver.Set_doorOpen_lamp(0)
		driver.Set_MotorDirection(elevator.CurrentDirection)

		if(elevator.CurrentDirection == DirStop){
			elevator.ElevatorState = idle
		}
		else {
			elevator.ElevatorState = moving
		}

		break
	default:
		break
	}
}