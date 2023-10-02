package main
 
import (
//        "fmt"
	"log"
	"strings"
	"os"
	"math/rand"
	"path/filepath"
	"time"
//	"strconv"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 10                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)
 
func main() {
	
     for i := 0; i < 99999; {
        pattern := "/etc/test*"
        files, err := filepath.Glob(pattern)
        if err != nil {
           log.Fatal(err)
           return
        }
        rand.Seed(time.Now().UnixNano())
        filefound := strings.Split(files[0]," ")
        oldLocation := filefound[0] //strings.Split(filefound[0]," ")
        newLocation := "/etc/test-"+RandStringBytesMask(rand.Intn(30 - 10) + 10)
        rand.Seed(time.Now().UnixNano())
        log.Println(rand.Intn(30 - 10) + 10)
        log.Println(oldLocation+" "+newLocation)
        errRename := os.Rename(oldLocation, newLocation)
        if errRename != nil {
                log.Fatal(errRename)
        }
	time.Sleep(3 * time.Second)

     }

}


func RandStringBytesMask(n int) string {
    b := make([]byte, n)
    for i := 0; i < n; {
        if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i++
        }
    }
    return string(b)
}
