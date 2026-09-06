package main

import (
	"bufio"
	// "fmt"
	"log"
	"os"
)


func main() {
	filename:="dirlist.txt"
	file,err := os.Open(filename)
	if err!=nil {
		// fmt.Errorf("failed to read %s : %w",filename,err)
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var directories []string
	for scanner.Scan() {
		directories=append(directories,scanner.Text())
	}
	
	// check if each directory exists
	for _,directory := range directories {
		_,err := os.Stat(directory)
		if err!=nil {
			log.Fatal(err)
		}
	}

}
