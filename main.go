package main

import (
	"fmt"
	dotsync "github.com/sugvn/dotsync/core"
	"flag"
)

func main() {
	filename := flag.String("file","dirlist.txt","specify the file that lists the directories to track")
	action := flag.String("action","none","specify whether to pull/push")
	flag.Parse()
	dir_map := dotsync.BuildDirMap(*filename)
	switch *action {
	case "none":
		fmt.Println("No operation provided")
	case "pull":
		err := dotsync.PullChanges(dir_map)
		dotsync.FatalOnErr(err)
	case "push":
		err := dotsync.PushChanges(dir_map)
		dotsync.FatalOnErr(err)
	default:
		fmt.Println("Invalid action")
	}
}
