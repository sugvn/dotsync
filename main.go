package main

import (
	"fmt"
	dotsync "github.com/sugvn/dotsync/core"
)

func main() {
	filename := "dirlist.txt"
	dir_map,err := dotsync.BuildDirMap(filename)
	if err != nil {
		fmt.Println("Error building directory map:",err.Error())
		return
	}
	// err = dotsync.PullChanges(dir_map)
	// if err != nil {
	// 	fmt.Println("Error Pulling Changes:")
	// 	fmt.Println(err.Error())
	// 	return
	// }
}
