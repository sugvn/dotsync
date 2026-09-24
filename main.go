package main

import (
	"fmt"
	dotsync "github.com/sugvn/dotsync/core"
	"flag"
)

func main() {
	filename := flag.String("file","dirlist.txt","specify the file that lists the directories to track")
	action := flag.String("action","none","specify whether to pull or push")
	flag.Parse()
	dir_map := dotsync.BuildDirMap(*filename)
	if *action=="none" {
		fmt.Println("No operation provided")
		return
	}
	if *action=="pull" {
		err := dotsync.PullChanges(dir_map)
		if err != nil {
			fmt.Println("Error Pulling Changes:")
			fmt.Println(err.Error())
			return
		}
	} else if *action=="push" {
		fmt.Println("Unimplemented")
	} else {
		fmt.Println("no action named ",*action)
	}
}
