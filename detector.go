package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
//        "syscall"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
        ps "github.com/mitchellh/go-ps"
        "log"
//        "github.com/fsnotify/fsevents"
)

//
var watcher *fsnotify.Watcher

// main
func main() {
	//load blacklist process
        blackbyte, errbl := ioutil.ReadFile("blacklist.json")
        if errbl != nil {
          fmt.Println(errbl)
        }
        var blacklist []string
	json.Unmarshal(blackbyte, &blacklist)
	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("/etc", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %s : %#v\n", event.Op, event)
			if event.Op&fsnotify.Rename == fsnotify.Rename {
                           fmt.Printf("<< WARNING KILL THIS!")
                           fmt.Printf("%+q", event)
   			   //pid := syscall.Getgid()
			   //fmt.Println("parent process ID:", pid)
			   for _, infect := range blacklist {
// 				panic(infect)
  				listproc(infect)
			   }
			}	
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func listproc(infect string){
     processList, err :=  ps.Processes()
     if err != nil {
    	log.Println("ps.Processes() failed, system not identified")
        return
     }

     for x := range processList {
	var process ps.Process
     	process = processList[x]
        inspect := strings.ToLower(strings.Trim(process.Executable(), " "))
        infectproc := strings.ToLower(strings.Trim(infect, " "))
//        log.Println(inspect+" -- "+infectproc)
        if (strings.Contains(inspect, infectproc)){
   	   log.Println("Dangerouse Process found : %d\t%s\n",process.Pid(),process.Executable())
        }
  	//log.Printf("%d\t%s\n",process.Pid(),process.Executable())
     }

}

