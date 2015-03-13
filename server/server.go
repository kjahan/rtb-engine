package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/bsm/openrtb"
    "github.com/kjahan/adengine/writer"
//    "time"
//    "math/rand"
    "os"
    "flag"
    "strings"
    "io"
    "encoding/json"
    "io/ioutil"
    "bytes"
    elastigo "github.com/mattbaird/elastigo/lib"
//    "strconv"
)

var (
    Info    *log.Logger
    Error   *log.Logger
    ESConn    *elastigo.Conn
    BidId   int
    Debug bool
)

type Bidval struct {
    Bid float32
}

func usage() {
    fmt.Fprintf(os.Stderr, "usage: %s port_no\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(2)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Init(
    infoHandle io.Writer,
    errorHandle io.Writer) {

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
    //turn off debug mode
    Debug = false
    //connect to db
    //ESConn = writer.ConnectES()
    writer.ConnectMongo()
}

func Cleanup() {
    writer.CloseMongo()
}

//send a request to bid optimizer service
func getBid(bid_req *openrtb.Request) float32{
    optimalAPI := "http://127.0.0.1:5001/bidders/const"	//const bidder
    //optimalAPI := "http://127.0.0.1:5001/bidders/lin"	//linear bidder
    js, err := json.Marshal(bid_req)
    check(err)
    //fmt.Println(js)
    payload := bytes.NewReader(js)
    resp, err := http.Post(optimalAPI, "application/json", payload)
    check(err)
    if resp.StatusCode != 200 {
        fmt.Println("Response with error!")
    }
    //read bid response body
    body, err := ioutil.ReadAll(resp.Body)
    check(err)
    //fmt.Println(string(body))
    res := &Bidval{}
    json.Unmarshal(body, &res)
    defer resp.Body.Close()
    return res.Bid
}

//index bid req in es
func WriteES(bid_req *openrtb.Request, id string) {
    writer.WriteES(ESConn, bid_req, id)
}

//rtb bid request handler
func BidRequestHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    //fmt.Println(r.Body)
    //amt := time.Duration(rand.Intn(120))
    //time.Sleep(time.Millisecond * amt)
    bid_req, err := openrtb.ParseRequest(r.Body)
    if err != nil {
        //Error.Println("ERROR %s", err.Error())
        panic(err)
    } else {
	//fmt.Println(bid_req)
        if Debug {
	    Info.Println("INFO  Received bid request", *bid_req.Id)
        }
    }
    //store bid req into data store
    //BidId += 1
    //fmt.Println("BidId=", BidId)
    writer.WriteMongo(bid_req)	//write bid reqs into mongo
    //writer.WriteCompressedMongo(bid_req)  //write compressed bid req into db
    //WriteES(bid_req, strconv.Itoa(BidId))	//write bid reqs into es
    //get optimal bid value
    opt_bidval := getBid(bid_req)
    if Debug {
        fmt.Println("optimal bid value:", opt_bidval)
    }
    //w.WriteHeader(204) // respond with 'no bid'
    //bid response
    var bid_resp *openrtb.Response
    bid_resp = new(openrtb.Response)
    bid_resp.Id = new(string)
    *bid_resp.Id = "1"
    bid := (&openrtb.Bid{}).SetID("BIDID").SetImpID("IMPID").SetPrice(opt_bidval)
    sb := openrtb.Seatbid{}
    sb.Bid = append(sb.Bid, *bid)
    bid_resp.Seatbid = append(bid_resp.Seatbid, sb)
    js, err := json.Marshal(bid_resp)
    if err != nil {
        Error.Println("ERROR %s %s", err.Error(), http.StatusInternalServerError)
    }
    if Debug {
        fmt.Println(string(js))
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func main(){
    if len(os.Args) != 2 {
	usage()
        os.Exit(1)
    }
    //setup logging
    info_file, err := os.OpenFile("info_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open info log file:", err)
    }
    err_file, err := os.OpenFile("err_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file:", err)
    }
    Init(info_file, err_file)
    //countResponse, _ := ESConn.Count("rtb", "bid_request", nil, nil)
    //fmt.Println("Count=", countResponse.Count)
    //BidId = countResponse.Count
    fmt.Printf("running http server on port %s\n", os.Args[1]);
    data := []string{":", os.Args[1]}
    port := strings.Join(data, "")
    //fmt.Printf(port);
    //register bid req handler --> endpoint: http://127.0.0.1/bid
    http.HandleFunc("/bid", BidRequestHandler)
    //run rtb server
    err = http.ListenAndServe(port, nil)
    if err != nil {
	log.Fatalln("ListenAndServe: ", err)
    }
    Cleanup()
}
