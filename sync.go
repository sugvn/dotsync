package dotsync

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	cp "github.com/otiai10/copy"
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

func printDirMap(dir_map map[string]string){
	for basename,path := range dir_map {
		fmt.Println(basename," : ",path)
	}
}

func cleanUp(path string) error{
	err:= os.RemoveAll(path)
	if err!=nil {
		fmt.Printf("cleanUp: Failed to remove directory: %s",path)
		return err
	}
	return nil
}

func revertFromTmp(path string,dir_map map[string]string) error {
	for basename := range dir_map {		tmpdir_basename:= path + "/" + basename
		err:= os.Rename(tmpdir_basename,basename)
		if err!=nil {
			fmt.Println("Failed to move ",tmpdir_basename," to ",basename)
			return err
		}
	}
	return nil
}

// create a tmp directory in the current directory
// for each item in the map,copy all files in the dotfiles/ into tmp directory
// continue if succeeds,else abort
// mv the files from the original location into the dotfiles directory
// continue if succeeds,else copy the files from tmp into dotfiles
// delete tmp
func PullChanges(dir_map map[string]string) error {

	// for key,val in dir_map: copy file(val) to file(key)	
	tmpdir,err := os.MkdirTemp("./","dotsync-tmp")
	if err!=nil {
		return fmt.Errorf("Failed to create tmp directory: %w",err)
	}
	fmt.Println("Temp Directory Name:",tmpdir)

	// mv all files into tmp
	for basename := range dir_map {		
		tmpdir_basename := tmpdir + "/" + basename
		err := os.Rename(basename,tmpdir_basename)
		if err!=nil {
			fmt.Printf("failed to move %s into tmp: %s \n",basename,err.Error())
		}
	}

	// copy all files from original path to root of dotfiles/
	for basename,path := range dir_map {
		err:= cp.Copy(path,basename)	
		if err!=nil {
			fmt.Printf("Failed to copy %s to %s: %s",path,basename,err.Error())
			err1:= revertFromTmp(tmpdir,dir_map)
			if err1!=nil {
				fmt.Println("reverting from tmp failed: ",err1.Error())
			} else {
				fmt.Println("Cleaning up tmp")
				cleanUp(tmpdir)
			}
			return err
		}
	}
	cleanUp(tmpdir)
	return nil
}

func PushChanges(dir_map map[string]string) error {
	//for key,val in dir_map: copy file(key) to file(val)
	return nil
}
