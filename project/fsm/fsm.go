package fsm

import (
	def "../definitions"
	"../driver"
	"../queue"
	"fmt"
	"time"
)

func FsmFloorArrival(new_floor int, elevator *def.Elevator) {
	if new_floor == -1 {
		fmt.Print("Run FSM_floor_arrival while not on floor\n")
	} else {
		//fmt.Print("FSM_floor_arrival\n")
		driver.SetFloorIndicator(new_floor)
		elevator.Last_floor = new_floor
		elevator.Motor_stop_timer.Stop()
		switch elevator.Elevator_state {
		case def.Moving:
			if queue.ShouldStop(*elevator) {
				driver.SetMotorDirection(def.Dir_stop)
				queue.ClearAtFloor(elevator, new_floor)
				driver.ClearLightsAtFloor(elevator.Last_floor)
				driver.SetDoorOpenLamp(1)
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

func FsmNextOrder(elevator *def.Elevator, next_order def.Order) { //arbitrator decides where we should go next
	fmt.Print("FSM_next_order\n")
	driver.SetButtonLamp(next_order, 1)

	switch elevator.Elevator_state {
	case def.Idle:
		queue.Enqueue(elevator, next_order)
		if next_order.Floor == elevator.Last_floor {
			queue.ClearAtFloor(elevator, elevator.Last_floor)
			driver.ClearLightsAtFloor(elevator.Last_floor)
			elevator.Door_timer.Reset(3 * time.Second)
			driver.SetDoorOpenLamp(1)
			elevator.Elevator_state = def.Stop_on_floor
		} else {
			if next_order.Floor > elevator.Last_floor {
				elevator.Current_direction = def.Dir_up
				driver.SetMotorDirection(elevator.Current_direction)
			} else {
				elevator.Current_direction = def.Dir_down
				driver.SetMotorDirection(elevator.Current_direction)
			}

		}
		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(4 * time.Second)
			fmt.Print("FSM_next_order: Reset motor_timer\n")
		}
	case def.Moving:
	case def.Stop_on_floor:
		queue.ClearAtFloor(elevator, elevator.Last_floor)
		driver.ClearLightsAtFloor(elevator.Last_floor)
		elevator.Door_timer.Reset(3 * time.Second)
	case def.Motor_stop:
		if next_order.Type == def.Buttoncall_internal {
			queue.Enqueue(elevator, next_order)
		}

	default:
		break
	}
	queue.PrintQueue(*elevator)
}

func FsmOnDoorTimeout(elevator *def.Elevator) {
	fmt.Print("FSM_on_door_timeout\n")
	queue.PrintQueue(*elevator)
	driver.SetDoorOpenLamp(0)
	switch elevator.Elevator_state {
	case def.Stop_on_floor:
		elevator.Current_direction = queue.ChooseDirection(*elevator)
		driver.SetMotorDirection(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(8 * time.Second)
			fmt.Print("FSM_where_to_next: Reset motor_timer\n")
		}
		break
	case def.Idle:
		elevator.Current_direction = queue.ChooseDirection(*elevator)
		driver.SetMotorDirection(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(8 * time.Second)
			fmt.Print("FSM_on_door_timeout: Reset motor_timer\n")
		}
		break
	default:
		break
	}
}

func FsmWhereToNext(elevator def.Elevator) {
	switch elevator.Elevator_state {
	case def.Stop_on_floor:
		elevator.Current_direction = queue.ChooseDirection(elevator)
		driver.SetMotorDirection(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(8 * time.Second)
			fmt.Print("FSM_where_to_next: Reset motor_timer\n")
		}
		break
	case def.Idle:
		elevator.Current_direction = queue.ChooseDirection(elevator)
		driver.SetMotorDirection(elevator.Current_direction)

		if elevator.Current_direction == def.Dir_stop {
			elevator.Elevator_state = def.Idle
		} else {
			elevator.Elevator_state = def.Moving
			elevator.Motor_stop_timer.Reset(8 * time.Second)
			fmt.Print("FSM_on_door_timeout: Reset motor_timer\n")
		}
		break
	default:
		break
	}
}

func FsmMotorStop(elevator *def.Elevator) def.Elevator {
	fmt.Print("FSM_motor_stop\n")
	elevator.Current_direction = def.Dir_stop
	driver.SetMotorDirection(def.Dir_stop)

	elev := driver.ElevInitFromBackup()
	return elev

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
