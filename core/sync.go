package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	cp "github.com/otiai10/copy"
)

func isAllExist(paths []string) (bool,error) {
	for _,path := range paths {
		_,err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return false,nil
			}
			return false,err
		}
	}
	return true,nil
}

func FatalOnErr(err error){
	if err!=nil {
		log.Fatal(err)
	}
}

func BuildDirMap(filename string) map[string]string {

	file,err := os.Open(filename)
	FatalOnErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dir_map:=make(map[string]string)
	count_map:=make(map[string]int)
	for scanner.Scan() {
		path:=scanner.Text()

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
	
	return dir_map
}

func printDirMap(dir_map map[string]string){
	for basename,path := range dir_map {
		fmt.Println(basename," : ",path)
	}
}

func cleanUp(path string) error{
	err:= os.RemoveAll(path)
	if err!=nil {
		fmt.Printf("cleanUp: Failed to remove directory: %s\n",path)
		return err
	}
	return nil
}

func revertFromTmp(path string,dir_map map[string]string) error {
	for basename := range dir_map {		
		tmpdirBasename := filepath.Join(path,basename)
		renameErr := os.Rename(tmpdirBasename,basename)
		if renameErr != nil {
			if os.IsNotExist(renameErr) {
				continue
			}
			return renameErr
		}
	}
	return nil
}

// check write permission
func canWrite(path string) (bool,error) {
	info,err := os.Stat(path) 
	if err != nil {
		return false,err
	}

	// if directory
	if info.IsDir() {
		tmpDir,mkdirErr := os.MkdirTemp(path,"dotsync-tmp-*")
		if mkdirErr != nil {
			if os.IsPermission(mkdirErr) {
				return false,nil
			}
			return false,mkdirErr
		}
		cleanupErr := cleanUp(tmpDir)
		if cleanupErr != nil {
			return true,fmt.Errorf("Error cleaning up %s: %w",tmpDir,cleanupErr)
		}
		return true,nil
	}
	
	// if file
	file,openErr := os.OpenFile(path,os.O_WRONLY,0)
	if openErr != nil {
		if os.IsPermission(openErr) {
			return false,nil
		}
		return false,openErr
	}
	file.Close()
	return true,nil
}


// create a tmp directory in the current directory
// for each item in the map,copy all files in the dotfiles/ into tmp directory
// mv the files from the original location into the dotfiles directory
// delete tmp
func PullChanges(dir_map map[string]string) error {

	var paths []string
	for _,path := range dir_map {
		paths = append(paths,path)
	}
	isExist,existErr := isAllExist(paths)
	if existErr != nil {
		return existErr
	}
	if isExist == false {
		fmt.Println("One of the paths mentioned does not exist")
		return nil
	}
	// create a temporary directory
	tmpdir,mkdirErr := os.MkdirTemp("./","dotsync-tmp-*")
	if mkdirErr != nil {
		return fmt.Errorf("create tmp dir failed: %w",mkdirErr)
	}
	fmt.Println("Temp Directory Name:",tmpdir)

	// mv all files into tmp
	for basename := range dir_map {		
		tmpdirBasename := filepath.Join(tmpdir,basename)
		renameErr := os.Rename(basename,tmpdirBasename)
		if renameErr == nil {
			continue
		}
		if os.IsNotExist(renameErr) {
			fmt.Println("rename ",basename," into ",tmpdirBasename,": ",basename," does not exist")
			continue
		}

		revertErr := revertFromTmp(tmpdir,dir_map)
		if revertErr != nil {
			return fmt.Errorf("move failed:%w; revert failed:%v",renameErr,revertErr)
		}
		cleanErr := cleanUp(tmpdir)
		if cleanErr != nil {
			return fmt.Errorf("move failed:%w; cleanUp failed:%v",renameErr,cleanErr)
		}
		return fmt.Errorf("move failed:%w",renameErr)
	}

	// copy all files from original path to root of dotfiles/
	for basename,path := range dir_map {
		copyErr := cp.Copy(path,basename)	
		if copyErr != nil {

			revertErr := revertFromTmp(tmpdir,dir_map)
			if revertErr != nil {
				return fmt.Errorf("copying failed:%w; reverting failed:%v",copyErr,revertErr)
			}

			cleanErr := cleanUp(tmpdir)
			if cleanErr != nil {
				return fmt.Errorf("copying failed:%w; cleanUp failed:%v",copyErr,cleanErr)
			}

			return fmt.Errorf("copying failed:%w",copyErr)
		}
	}

	// clean up
	cleanErr := cleanUp(tmpdir)
	if cleanErr != nil {
		return fmt.Errorf("cleanUp failed:%w",cleanErr)
	}
	return nil
}

// for simplicity and clean transition of states,
// we can only push if there are basename directories of all entries specified in the list provided
// by the file
// i.e every entry in the file list must already be pulled before
// dont push partial list of entries
func PushChanges(dir_map map[string]string) error {

	// existence check
	var basenames []string
	for basename := range dir_map {
		basenames = append(basenames,basename)
	}
	isExist,existErr := isAllExist(basenames)
	if existErr != nil {
		return existErr
	}
	if isExist == false {
		fmt.Println("One of the paths mentioned does not exist")
		return nil
	}

	for basename,path := range dir_map {

		// check if writable,if true continue else skip 
		canWrite,writeErr := canWrite(path)
		fmt.Println("Can write to ",path,": ",canWrite)
		if !canWrite {
			if writeErr != nil {
				fmt.Printf("PushChanges: Error checking write permission for %s (skipping): %s\n",path,writeErr.Error())
			} else {
				fmt.Println("Not enough permission to write to ",path,": skipping")
			}
			continue
		}

		// remove <path>
		cleanupErr := cleanUp(path)
		if cleanupErr != nil {
			fmt.Printf("PushChanges: Error removing %s: %s ",path,cleanupErr.Error())
			return cleanupErr
		}

		// copy <basename> to <path>
		copyErr := cp.Copy(basename,path)
		if copyErr != nil {
			fmt.Println("failed copying ",basename," to ",path,": ",copyErr.Error())
			return copyErr
		}
		fmt.Println("copied ",basename,"to ",path)
	}
	
	return nil
}
