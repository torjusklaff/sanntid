package backup

import (
	def "../definitions"
	"log"
	"os"
	"strconv"
	"strings"
)

const filename = "log.txt"

func checkError(err error) {
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
}

func ToBackup(str string) {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println(str)
}

func ReadLastLine(string_size int64) string {
	string_size += 1
	f, err := os.Open(filename)
	checkError(err)
	defer f.Close()

	buf := make([]byte, string_size)
	stat, err := os.Stat(filename)
	start := stat.Size() - string_size
	n, err := f.ReadAt(buf, start)
	checkError(err)
	buf = buf[:n]
	return string(buf)
}

func BackupInternalQueue(elevator def.Elevator) {
	queue_string := QueueToString(elevator)
	ToBackup(queue_string)
}

func QueueToString(e def.Elevator) string {
	var queue_string string
	var order_string string
	for f := 0; f < def.N_floors; f++ {
		for btn := 0; btn < def.N_buttons; btn++ {
			if e.Queue[f][btn] == 1 {
				order_string = "1 "
			} else {
				order_string = "0 "
			}
			queue_string += order_string
		}
	}
	return queue_string
}

func QueueFromString(queue_string string) [4][3]int {
	queue := [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	index := 0

	queue_temp := strings.Fields(queue_string)

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			new_queue, _ := strconv.Atoi(queue_temp[index])
			queue[i][j] = new_queue
			index += 1
		}
	}
	return queue
}
