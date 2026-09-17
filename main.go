package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func handle(err error){
	if err!=nil {
		log.Fatal(err)
	}
}

func main() {
	filename:="dirlist.txt"
	file,err := os.Open(filename)
	handle(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dir_map:=make(map[string]string)
	count_map:=make(map[string]int)
	for scanner.Scan() {
		path:=scanner.Text()

		// existence check
		_,err := os.Stat(path)
		handle(err)

		basename := filepath.Base(path)

		if dir_map[basename]!="" {
			new_basename := basename + "_" + strconv.Itoa(count_map[basename])
			count_map[basename]+=1
			
			dir_map[new_basename]=path
		} else {
			dir_map[basename]=path
			count_map[basename]=1
		}
	}
	
	for basename,path := range dir_map {
		fmt.Println(basename," : ",path)
	}

}
