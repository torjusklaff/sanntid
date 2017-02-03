// Test-main for driver-files
import (
	"/driver"
)


func main(){
	driver.Elev_init()
	/*for {
		for floor := 0; floor < driver.N_floors; floor++ {
			for button := 0; button < driver.N_buttons; button++ {
				if driver.Get_button_signal(button, floor) {
					driver.Set_button_lamp(button, floor, 1)
				}
			}
		}
	}*/
	driver.Set_button_lamp(driver.Buttoncall_up, 1, 1)
}
