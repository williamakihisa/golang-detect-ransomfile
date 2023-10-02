package main

import (
    "log"
    "net/http"
    "sync"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "time"
    "strings"
    "strconv"
)

var conflist []string
var blockstatus int = 0
var counterblock int = 10

type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

func initConfig(){
   confbyte,err := ioutil.ReadFile("config.json")
   if err != nil {
	fmt.Println(err)
   }
   json.Unmarshal(confbyte,&conflist)
}

func clearBlockedIPs() {
    for {
        c = SafeCounter{v: make(map[string]int)}
        time.Sleep(time.Second * 5)
    }
}

func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()
	return c.v[key]
}

var c = SafeCounter{v: make(map[string]int)}

func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mux.Unlock()
}

func middleware(next http.Handler) http.Handler {
    initConfig()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddr := strings.Split(r.RemoteAddr, ":")[0]

		c.Inc(ipAddr)

        if c.Value(ipAddr) >= counterblock {
            http.Error(w, "", http.StatusTooManyRequests)
            return
		}
		if next == nil {
            http.DefaultServeMux.ServeHTTP(w, r)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main(){
  go clearBlockedIPs()
  http.HandleFunc("/saveES", poolHandler)
  if (blockstatus == 1){
    http.ListenAndServe(":9293", middleware(nil))
  }else{
    http.ListenAndServe(":9293", nil)
  }
}

func poolHandler(w http.ResponseWriter, r *http.Request) {
    host, cekhost := r.URL.Query()["host"]
    event, cekevent := r.URL.Query()["event"]
    ipAddr := strings.Split(r.RemoteAddr, ":")[0]
    var hostreq string = ""
    var eventreq string = ""
    now := time.Now()
    current := time.Now().Round(0)
    curdatetime := strings.Replace(current.Format("2006-01-02 15:04:05 GMT+07")," GMT+07", "", 1)
    timestamp := now.Unix()
    idtimestamp := strconv.FormatInt(timestamp, 10)
    if cekhost {hostreq = host[0]}
    if cekevent {eventreq = event[0]}

    var reqBody string = `{
       "id": `+idtimestamp+`,
       "server_host": "`+hostreq+`",
       "ip_address" : "`+ipAddr+`",
       "event" : "`+eventreq+`",
       "datetime": "`+curdatetime+`"
    }`

    urlelastic := "http://YOURELASTIC:9200/YOURINDEX/_doc/"+idtimestamp /// << change Elasticsearch Endpoint and Index
    respSave := callElastic(urlelastic,reqBody,"PUT")
    fmt.Fprint(w, respSave)
}


func callElastic(url string, reqBody string, method string) string{
   var resp string = ""
         payload := strings.NewReader(reqBody)
   	 client := &http.Client {
	 }
	 req, err := http.NewRequest(method, url, payload)

	 if err != nil {
	 log.Fatal(err)
	 }
	 req.Header.Add("Content-Type", "application/json")
         req.Header.Add("Authorization", "Basic bXBpOjNzNHAx")
	 res, err := client.Do(req)
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer res.Body.Close()

	 body, err := ioutil.ReadAll(res.Body)
	 if err != nil {
		 log.Fatal(err)
	 }else{
 		resp = string(body)
	 }
  return resp
}
