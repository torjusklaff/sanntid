package backup

import (
	"log"
	"os"
)

const filename = "log.txt"
const string_size = 62 		// Leser fra og med size minus string_size (aka siste linje)

func check_error(err error){
	if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
}


func To_backup(str string) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check_error(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println(str)
}


func Read_last_line() line string {
	f, err := os.Open(filename)
	check_error(err)
	defer f.Close()

	buf := make([]byte, string_size)
	stat, err := os.Stat(filename)
	start := stat.Size() - string_size
	line, err = f.ReadAt(buf, start)
	check_error(err)
}