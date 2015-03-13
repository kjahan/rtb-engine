package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "adengine/openrtb"
//    "time"
)

type Bidval struct {
    Bid float32
}

var (
    Debug bool
)

func ComputeBid(bid_req *openrtb.Request) *Bidval {
    b, _ := json.Marshal(bid_req)
    if Debug {
        fmt.Println(fmt.Sprintf("%#v", bid_req))
        fmt.Println(string(b))
        fmt.Println("INFO: bid request id =", *bid_req.Id)
        if bid_req.Device != nil && bid_req.Device.Make != nil {
            fmt.Println("INFO: dev make =", *bid_req.Device.Make)
        } else {
            fmt.Println("INFO: dev make is null!")
        }
        if bid_req.Device != nil && bid_req.Device.Model != nil {
            fmt.Println("INFO: dev model =", *bid_req.Device.Model)
        }
        if bid_req.Device != nil && bid_req.Device.Os != nil {
            fmt.Println("INFO: dev OS =", *bid_req.Device.Os)
        }
        if bid_req.Device != nil && bid_req.Device.Osv != nil {
            fmt.Println("INFO: dev OS version =", *bid_req.Device.Osv)
        }
        if bid_req.Device != nil && bid_req.Device.Ua != nil {
            fmt.Println("INFO: user agent =", *bid_req.Device.Ua)
        }
        if bid_req.Device != nil && bid_req.Device.Language != nil {
            fmt.Println("INFO: lang =", *bid_req.Device.Language)
        }
        if bid_req.Device != nil && bid_req.Device.Connectiontype != nil {
            fmt.Println("INFO: conn type =", *bid_req.Device.Connectiontype)
        }
        if bid_req.Device != nil && bid_req.Device.Devicetype != nil {
            fmt.Println("INFO: dev type =", *bid_req.Device.Devicetype)
        }
        if bid_req.Device != nil && bid_req.Device.Js != nil {
            fmt.Println("INFO: js =", *bid_req.Device.Js)
        }
        if bid_req.Device != nil && bid_req.Device.Carrier != nil {
            fmt.Println("INFO: carrier =", *bid_req.Device.Carrier)
        }
        if bid_req.Device != nil && bid_req.Device.Geo.Country != nil {
            fmt.Println("INFO: country =", *bid_req.Device.Geo.Country)
        }
        if bid_req.Device != nil && bid_req.Device.Geo.Region != nil {
            fmt.Println("INFO: region =", *bid_req.Device.Geo.Region)
        }
        if bid_req.App != nil && bid_req.App.Id != nil {
            fmt.Println("INFO: App id =", *bid_req.App.Id)
        }
        if bid_req.App != nil && bid_req.App.Name != nil {
            fmt.Println("INFO: App name =", *bid_req.App.Name)
        }
        if bid_req.App != nil && bid_req.App.Domain != nil{
            fmt.Println("INFO: App domain =", *bid_req.App.Domain)
        }
        if bid_req.User != nil && bid_req.User.Id != nil{
            fmt.Println("INFO: user id =", *bid_req.User.Id)
        }
    }
    //say there is 75ms delay for computing optimal bid in the worst case!
    //amt := time.Duration(75 * time.Millisecond)
    //time.Sleep(amt)
    bid_val := new(Bidval)
    bid_val.Bid = float32(0.5)	//here we make a call to bid optimizer layer!
    return bid_val
}

//optimizer web server
func BidOptimizer(w http.ResponseWriter, r *http.Request) {
    bid_req, err := openrtb.ParseRequest(r.Body)
    if err != nil {
        fmt.Println("ERROR", err.Error())
    }
    bid_val := ComputeBid(bid_req)
    js, err := json.Marshal(bid_val)
    if err != nil {
        log.Fatal("InternalServerError: ", err)
    }
    //fmt.Println(string(js))
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func main(){
    http.HandleFunc("/optimizer", BidOptimizer)
    err := http.ListenAndServe(":5001", nil)
    Debug = false
    if err != nil {
	log.Fatal("ListenAndServe: ", err)
    }
}
