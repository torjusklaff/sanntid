import (
	"../driver"
)

type elev_states int

const (
	idle elev_states = iota
	door_open
	moving
)


type Elevator struct {
	last_floor int
	current_direction motor_direction
	queue int
	elevator_state elev_states
	door_open_duration float
}

var elev_data elevator_data


func get_elev_data(){
	return elev_data
}

func set_elev_data(floor last_floor, direction motor_direction){
	elev_data{last_floor: floor}
	elev_data{motor_direction: direction}
}


func init() {
	elev_init()
	set_elev_data(get_floor_sensor_signal(), dir_stop)
}
