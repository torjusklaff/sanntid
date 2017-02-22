package definitions

const (
	N_floors  = int(C.N_FLOORS)
	N_buttons = int(C.N_BUTTONS)
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
	Type  Button_type
	Floor int
}