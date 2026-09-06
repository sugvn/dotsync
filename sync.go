package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	var paths []string
	for scanner.Scan() {
		paths=append(paths,scanner.Text())
	}
	
	// check if each directory exists
	for _,path := range paths {
		_,err := os.Stat(path)
		if err!=nil {
			log.Fatal(err)
		}
	}

	var basenames []string
	for _,path := range paths {
		basename := filepath.Base(path)
		basenames = append(basenames,basename)	
	}

	for _,basename := range basenames {
		fmt.Println(basename)
	}


}
