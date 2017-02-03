// Test-main for driver-files
import (
	"/driver"
)


func main(){
	elev_init()

	for floor := 0; floor < n_floors; floor++ {
		for button := 0; button < n_button; button++ {
			if get_button_signal(button, floor){
				set_button_lamp(button, floor, 1)
			}
		}
}
