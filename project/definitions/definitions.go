package definitions

const (
	N_floors    = 4
	N_buttons   = 3
	N_elevators = 3
)

type Motor_direction int

const (
	Dir_down Motor_direction = -1
	Dir_stop Motor_direction = 0
	Dir_up   Motor_direction = 1
)

type Button_type int

const (
	Buttoncall_down     Button_type = 1
	Buttoncall_up       Button_type = 0
	Buttoncall_internal Button_type = 2
)

type Order_button struct {
	Type     Button_type
	Floor    int
	Internal bool
	Id       string
}

type Elev_states int

const (
	Idle      Elev_states = iota
	Door_open             //Not really necessary, look into it (Change to On_floor)
	Moving
)

type Elevator struct {
	Last_floor        int
	Current_direction Motor_direction
	Queue             [N_floors][N_buttons]int
	Elevator_state    Elev_states
	Id                string
}

type Cost struct {
	Cost          float64
	Current_order Order_button
	Id            string
}
