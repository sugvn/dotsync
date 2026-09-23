package dotsync

import "testing"

func TestBuildDirMap(t *testing.T){
	filename:="dirlist.txt"
	dir_map,err := BuildDirMap(filename)
	Handle(err)
	printDirMap(dir_map)
}

func TestPullChanges(t *testing.T){
	filename:= "dirlist.txt"
	dir_map,err:= BuildDirMap(filename)
	Handle(err)
	err1:= PullChanges(dir_map)
	Handle(err1)
}
