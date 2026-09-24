package core

import "testing"

func TestBuildDirMap(t *testing.T){
	filename:="../dirlist.txt"
	dir_map := BuildDirMap(filename)
	printDirMap(dir_map)
}

func TestPullChanges(t *testing.T){
	filename:= "../dirlist.txt"
	dir_map := BuildDirMap(filename)
	err := PullChanges(dir_map)
	FatalOnErr(err)
}
