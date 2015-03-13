package main

import (
    "fmt"
    "net/http"
    "bytes"
    "math/rand"
    "os"
    "flag"
    "strings"
    "io/ioutil"
    "bufio"
    "io"
    "time"
    "sync"
)

var (
    Debug bool
    wg sync.WaitGroup
)

func usage() {
    fmt.Fprintf(os.Stderr, "usage: %s port\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(2)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func timeTrack(start time.Time, name string, messages chan string) {
    elapsed := time.Since(start)
    //fmt.Println(name, " took ", elapsed)
    messages <- fmt.Sprintf("%s took %v", name, elapsed)
}

func getAppBidReq(filename string) []byte {
    // read the whole file
    payload, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    return payload
}

func getSimpleBanner(filename string) []byte {
    // read whole the file
    payload, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    //let's pick a random bid id
    sampleIds := [][]byte{[]byte(`9994534625259`), []byte(`1234534625000`), []byte(`5674534625890`)}
    inx := rand.Intn(len(sampleIds))
    copy(payload[11:24], sampleIds[inx])
    //s = string(payload[11:24])
    //fmt.Println(s)
    return payload
}

func sendBidReq(BidAPI string, file_name string, messages chan string) {
    defer wg.Done()
    defer timeTrack(time.Now(), "sendBidReq", messages)
    writes := 0
    fp, err := os.Open(file_name)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer fp.Close()
    r := bufio.NewReader(fp)
    //read payload
    //payload := getAppBidReq(file_name)
    for {
        //s, err :=  r.ReadString('\n')
        payload, err := r.ReadBytes(0x0A)
        if err == io.EOF {
	    break
	}
        //payload :=  []byte(s[:len(s) - 1])
	if err != nil {
	    panic("GetLines: " + err.Error())
        }
/*
        if Debug {
            fmt.Println("payload=", s[:len(s)-1])      //bid req payload
        }
*/
        //buf := bytes.NewReader(payload)
        buf := bytes.NewReader(payload[:len(payload) - 1])
        resp, err := http.Post(BidAPI, "application/json", buf)
        check(err)
        if resp.StatusCode == 204 {
            fmt.Println("Response with no bid!")
        }else{
            //read bid response body
            body, err := ioutil.ReadAll(resp.Body)
            check(err)
            if Debug {
                fmt.Println(string(body))
            }
        }
        defer resp.Body.Close()
        writes += 1
        if writes >= 50000 {
            break
        }
    }
}

//web client: SSP
/*
func sendRTBRequest(no int, BidAPI string, file_name string) {
    for i := 0; i < no - 1; i++ {
        go sendBidReq(BidAPI, file_name)
    }
    sendBidReq(BidAPI, file_name)
}
*/

func main() {
    if len(os.Args) != 2 {
        usage()
        os.Exit(1)
    }
    fmt.Printf("running http clients on port %s\n", os.Args[1]);
    messages := make(chan string)
    //we have starting 5 goroutines
    wg.Add(5)
    //req_no, err := strconv.Atoi(os.Args[1])
    //check(err)
    file_name_1 := "/opt/data/train_1458_logs_1.txt"
    file_name_2 := "/opt/data/train_1458_logs_2.txt"
    file_name_3 := "/opt/data/train_1458_logs_3.txt"
    file_name_4 := "/opt/data/train_1458_logs_4.txt"
    file_name_5 := "/opt/data/train_1458_logs_5.txt"
    data := []string{"http://52.10.197.103:", os.Args[1], "/bid"}
    BidAPI := strings.Join(data, "")
    //turn off debug mode
    Debug = false
    //run two go routines to send bid requests
    go sendBidReq(BidAPI, file_name_1, messages)
    go sendBidReq(BidAPI, file_name_2, messages)
    go sendBidReq(BidAPI, file_name_3, messages)
    go sendBidReq(BidAPI, file_name_4, messages)
    go sendBidReq(BidAPI, file_name_5, messages)
    go func() {
        for i := range messages {
            fmt.Println(i)
        }
    }()
    wg.Wait()
}
