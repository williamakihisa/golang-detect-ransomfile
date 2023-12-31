package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
  ps "github.com/mitchellh/go-ps"
  "log"
	"strconv"
	"net/http"
	"time"
	"os/exec"
)
//setting sections <- change this
var mailnotif string = "YOURMAILADDRESS@MAIL.COM"
var hostserve string = "SERVER_NAME"
var mailserve string = "MAIL_SERVER"
var mailkey string = "MAIL_API_KEY"
var masterAPI string = "MASTER_PROTECT_POOLER_SERVICE"
//END setting sections

var watcher *fsnotify.Watcher
var blacklist []string

var event_notif string


// main
func main() {
	//load blacklist process
        blackbyte, errbl := ioutil.ReadFile("blacklist.json")
        whitebyte,errwh := ioutil.ReadFile("whitelist.json")

        if errbl != nil {
          fmt.Println(errbl)
        }
        var whitelist []string
        if errwh != nil {
	  fmt.Println(errwh)
        }else{
          json.Unmarshal(whitebyte, &whitelist)
          if (len(whitelist) == 0) {
             listproc("whitelisitignorexyz", 1)
             whitebyte,errwh := ioutil.ReadFile("whitelist.json")
             if errwh != nil {
               fmt.Println(errwh)
             }
	     json.Unmarshal(whitebyte, &whitelist)
          }
        }
	json.Unmarshal(blackbyte, &blacklist)
	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("/etc", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}


	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %s : %#v\n", event.Op, event)
			if event.Op&fsnotify.Rename == fsnotify.Rename {
                           fmt.Printf("<< ALERT SUSPICIOUS PROCESS!!")
                           fmt.Println("%+q", event)
			   pidInfect := 0
			   for _, infect := range blacklist {
  				pidInfect = listproc(infect, 0)
			   }
			   if (pidInfect == 0){
				log.Println("no infected blacklist process found, begin checking whitelist process...")
				log.Println("%+q",whitelist)
				var diffres []int
				diffres = getDiffProcess(whitelist)
                                if (len(diffres) > 0) {
				 log.Println("found alien process > %+q",diffres)

			         for _, killID := range diffres {
			 	   proc, errproc := os.FindProcess(killID)
        			   if errproc != nil {
	   			     fmt.Println("error process find :", errproc)
				   }
					 printPath(killID)
				   proc.Kill()
     				 }
				 //end kill process
				}

			   }else{
				log.Println(" process id  = "+strconv.Itoa(pidInfect))
				//kill process
				proc, errproc := os.FindProcess(pidInfect)
        			if errproc != nil {
           			   fmt.Println("error process find :", errproc)
        			}
							printPath(pidInfect)
        			proc.Kill()
			   }
        fmt.Println("found virus time to notif > ",event_notif)
			  //submit mail
 			  datamail := map[string]string{}

			  datamail["to"] = mailnotif
			  datamail["subject"] = "Suspicion Process Found And Cleared in : "+hostserve
		          datamail["html"] = "<p><b>"+event_notif+"</b></p>"
		          datamail["company"] = "MNC Portal Indonesia"
			  datamail["sendername"] = "MNC Server System"
			  jsonMail, errjsmail := json.Marshal(datamail)
			  if errjsmail != nil {
			     fmt.Println("ERR MAIL : ",errjsmail)
			  }
	  		  fmt.Println(string(jsonMail))

  url := mailserve
  method := "POST"
  payload := strings.NewReader(string(jsonMail))
  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("x-apikey", mailkey)
  req.Header.Add("Content-Type", " application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
  updateDetection(event_notif)

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

func listproc(infect string, typeproc int) int{

     pid := 0
     processList, err :=  ps.Processes()
     if err != nil {
    	log.Println("ps.Processes() failed, system not identified")
        return 0
     }

     var tmpWhite []string
     y := 0

     for x := range processList {
			var process ps.Process
     	process = processList[x]
        inspect := strings.ToLower(strings.Trim(process.Executable(), " "))
        infectproc := strings.ToLower(strings.Trim(infect, " "))
        if (strings.Contains(inspect, infectproc)){
           event_notif = event_notif + "-(" +inspect+":"+strconv.Itoa(process.Pid())+")"
           if typeproc == 0 {
	     return process.Pid()
           }
        }else{
	  		  tmpWhite = append(tmpWhite, inspect)
          y++
	      }
     }
   if typeproc != 0 {
      writeJSONToken(tmpWhite,"whitelist.json")
   }
   return pid
}


func writeJSONToken(list []string, filename string){
  jsonString, _ := json.Marshal(list)
  ioutil.WriteFile(filename, jsonString, os.ModePerm)
}

func getDiffProcess(whitelist []string) []int{
     var diffArr []int
     processList, err :=  ps.Processes()
     if err != nil {
        log.Println("ps.Processes() failed, system not identified")
     }
     for x := range processList {
	var process ps.Process
        safestatus := 0
     	process = processList[x]
        inspect := strings.ToLower(strings.Trim(process.Executable(), " "))
        for _, safeproc := range whitelist {
            if (safeproc == inspect) {
		safestatus = 1
	    }else{
              //check if system worker multiplied
	      if strings.Contains(inspect, safeproc) {
		safestatus = 1
              }
            }
	}
	if (safestatus == 0) {
           updateBlacklist(inspect)
	   diffArr = append(diffArr, process.Pid())
	   event_notif = event_notif + "-(" + inspect + ":" + strconv.Itoa(process.Pid()) + ")"
        }
     }
     return diffArr
}

func updateBlacklist(procname string){
  blacklist = append(blacklist,procname)
  writeJSONToken(blacklist,"blacklist.json")
}

func updateDetection(eventDetail string){

  var detectlist []string
  detectbyte, errdet := ioutil.ReadFile("detectlist.json")
  if errdet != nil {
     fmt.Println(errdet)
  }
  json.Unmarshal(detectbyte, &detectlist)
    now := time.Now().Round(0)

    t := now.Format("2006-01-02 15:04:05 GMT+07")
    itemdetect := t+"="+eventDetail
    detectlist = append(detectlist,itemdetect)
    writeJSONToken(detectlist, "detectlist.json")
   //send detection notif to master protector API
   url := masterAPI+"?host="+hostserve+"&event="+eventDetail
   method := "GET"
   payload := strings.NewReader("")
   client := &http.Client {
   }
   req, err := http.NewRequest(method, url, payload)

   if err != nil {
      fmt.Println(err)
      return
   }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()
}

func printPath(pid int) {
    stringcommand := "/proc/"+strconv.Itoa(pid)+"/exe"

    prg := "ls"

    arg1 := "-l"
    arg2 := stringcommand

    cmd := exec.Command(prg, arg1, arg2)
    stdout, err := cmd.Output()

    if err != nil {
        log.Println(err.Error())
    }

    output := strings.Split(string(stdout)," -> ")
    infectpath := strings.Trim(output[1], " ")
    tmpinfect := strings.Split(infectpath,"/")
    filenameinfect := tmpinfect[(len(tmpinfect)-1)]
    cmdpwd := exec.Command("pwd")
    stdoutpwd, _ := cmdpwd.Output()
    qtine := fmt.Sprintf("%s%s%s%s",strings.TrimSuffix(string(stdoutpwd), "\n"),"/quarantine/",strings.TrimSuffix(filenameinfect, "\n"),".infect")
    log.Println(qtine)
    qtine = strings.TrimSuffix(qtine, "\n")
    errMove := os.Rename(strings.TrimSuffix(infectpath, "\n"), qtine)
    if errMove != nil {
       log.Println("failed to quarantine", errMove)
    }
}
