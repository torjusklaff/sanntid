package backup

import (
	"log"
	"os"
	def "../definitions"
	"strings"
	"strconv"
)

const filename = "log.txt"

func check_error(err error){
	if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
}


func To_backup(str string) {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check_error(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println(str)
}


func Read_last_line(string_size int64) string {
	string_size += 1
	f, err := os.Open(filename)
	check_error(err)
	defer f.Close()

	buf := make([]byte, string_size)
	stat, err := os.Stat(filename)
	start := stat.Size() - string_size
	n, err := f.ReadAt(buf, start)
	check_error(err)
	buf = buf[:n]
	return string(buf)
}

func Backup_internal_queue(elevator def.Elevator){
	queue_string := Queue_to_string(elevator)
	To_backup(queue_string)
}


func Queue_to_string(e def.Elevator) string {
	var queue_string string
	var order_string string
	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {
			if e.Queue[f][btn] == 1{
				order_string = "1 "
			} else {
				order_string = "0 "
			}
			queue_string += order_string
		}
	}
	return queue_string
}

func Queue_from_string(queue_string string) [4][3]int {
	queue := [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	index := 0

	queue_temp := strings.Fields(queue_string)
	
	for i:=0; i<4; i++{
		for j:=0; j<3; j++{
			new_queue, _ := strconv.Atoi(queue_temp[index])	
			queue[i][j] = new_queue
			index += 1
		}
	}
	return queue
}
