package elevator

import (
	"../driver"
)

type Elev_states int

const (
	idle Elev_states = iota
	door_open
	moving
)

type Elevator struct {
	last_floor         int
	current_direction  driver.Motor_direction
	queue              [driver.N_floors]int
	elevator_state     Elev_states
	door_open_duration float64
}

var elev_data Elevator

func Get_elev_data() Elevator {
	return elev_data
}

func Set_elev_data(floor int, direction driver.Motor_direction) {
	elev_data{last_floor: floor}
	elev_data{motor_direction: direction}
}

func Init() {
	elev_init()
	set_elev_data(get_floor_sensor_signal(), dir_stop)
}

