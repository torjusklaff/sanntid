import (
	"io"
	"log"
	"os"
)

// help function, checks errors
func check_error(err error){
	if err != nil {
        log.Fatal(err)
    }
}

// creates backup.txt if it doesnt exist
func new_backup_file(){
	if _, err := os.Stat("backup.txt"); os.IsNotExist(err) {
		newFile, err = os.Create("backup.txt")
		check_error(err)
		log.Println(newFile)
		newFile.Close()
	}
}

func read_backup(length int){
	file, err := os.Open("backup.txt")
	check_error(err)

	byte_slice := make([]byte, length)
	bytes_read, err := io.ReadFull(file, byte_slice)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
	return bytes_read
}


func read_full_backup(){
	file, err := os.Open("backup.txt")
	check_error(err)
	
	data, err := ioutil.ReadAll(file)
	check_error(err)
	
	file.Close()
	return data
}


func write_to_backup(to_file string) {
	file, err := os.OpenFile("backup.txt", os.O_RDWR|os.O_APPEND, 0660)
	check_error(err)
	defer file.Close()

	byte_slice := []byte(to_file)
	bytes_written, err := file.Write(byte_slice)
	check_error(err)
	file.Close()
}
