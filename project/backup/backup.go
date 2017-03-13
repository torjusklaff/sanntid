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
	f, err := os.OpenFile("log.txt", os.ORDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println(str)
}

func ReadLastLine(stringSize int64) string {
	stringSize += 1
	f, err := os.Open(filename)
	checkError(err)
	defer f.Close()

	buf := make([]byte, stringSize)
	stat, err := os.Stat(filename)
	start := stat.Size() - stringSize
	n, err := f.ReadAt(buf, start)
	checkError(err)
	buf = buf[:n]
	return string(buf)
}

func BackupInternalQueue(elevator def.Elevator) {
	queueString := QueueToString(elevator)
	ToBackup(queueString)
}

func QueueToString(e def.Elevator) string {
	var queueString string
	var orderString string
	for f := 0; f < def.NumFloors; f++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if e.Queue[f][btn] == 1 {
				orderString = "1 "
			} else {
				orderString = "0 "
			}
			queueString += orderString
		}
	}
	return queueString
}

func QueueFromString(queueString string) [4][3]int {
	queue := [4][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	index := 0

	tempQueue := strings.Fields(queueString)

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			nextInQueue, _ := strconv.Atoi(tempQueue[index])
			queue[i][j] = nextInQueue
			index += 1
		}
	}
	return queue
}
