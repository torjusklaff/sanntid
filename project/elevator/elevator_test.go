import (
	"../driver"
)

const deadline_period = time.Duration(5*driver.N_floors)*time.Second
const door_period = 3*time.Second

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


func Get_elev_data(){
	return elev_data
}

func Set_elev_data(floor last_floor, direction motor_direction){
	elev_data{last_floor: floor}
	elev_data{motor_direction: direction}
}


func Init(
	completed_floor chan<- int,
	missed_deadline chan<- bool,
	floor_reached <- chan int,
	new_target_floor <- chan int) {

	deadline_timer := time.NewTimer(deadline_period)
	deadline-timer.Stop()
	door_timer := time.NewTimer(door_period)
	door_timer.Stop()

	state := idle
	target_floor := -1

	for {
		
	}


}

