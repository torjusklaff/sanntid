import (
	"../driver"
)

type elevator_data struct {
	last_floor int
	current_direction motor_direction
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
