package backup

import (
	"log"
	"os"
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
