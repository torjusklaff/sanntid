import (
	"../driver"
)

const deadline_period = time.Duration(5*driver.NumFloors)*time.Second
const door_period = 3*time.Second

type elevatorStates int

const (
	idle elevatorStates = iota
	doorOpen
	moving
)


type Elevator struct {
	LastFloor int
	CurrentDirection MotorDirection
	queue int
	ElevatorState elevatorStates
	doorOpen_duration float
}

var elev_data elevator_data


func Get_elev_data(){
	return elev_data
}

func Set_elev_data(floor LastFloor, direction MotorDirection){
	elev_data{LastFloor: floor}
	elev_data{MotorDirection: direction}
}


func Init(
	completedFloor chan<- int,
	missed_deadline chan<- bool,
	floor_reached <- chan int,
	new_targetFloor <- chan int) {

	deadline_timer := time.NewTimer(deadline_period)
	deadline-timer.Stop()
	DoorTimer := time.NewTimer(door_period)
	DoorTimer.Stop()

	state := idle
	targetFloor := -1

	for {
		
	}
}

