package dotsync

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func Handle(err error){
	if err!=nil {
		log.Fatal(err)
	}
}

func BuildDirMap(filename string) (map[string]string,error) {

	file,err := os.Open(filename)
	Handle(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dir_map:=make(map[string]string)
	count_map:=make(map[string]int)
	for scanner.Scan() {
		path:=scanner.Text()

		// existence check
		_,err := os.Stat(path)
		Handle(err)

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
	
	return dir_map,nil
}

func PrintDirMap(dir_map map[string]string){
	for basename,path := range dir_map {
		fmt.Println(basename," : ",path)
	}
}
