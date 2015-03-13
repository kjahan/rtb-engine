package writer

import (
    "log"
    "gopkg.in/mgo.v2"
    "github.com/bsm/openrtb"
    "fmt"
    "bytes"
    "compress/zlib"
    "encoding/json"
)

type Person struct {
        Name string
        Phone string
}

var (
    MongoSession *mgo.Session
    MongoC  *mgo.Collection
)

func ConnectMongo() {
    //session, err := mgo.Dial("localhost")
    session, err := mgo.Dial("mongodb://localhost:27017")
    MongoSession = session
    if err != nil {
        panic(err)
    }
    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)
    MongoC = session.DB("rtb").C("bid_requests")
    //MongoC = session.DB("test").C("people")
    fmt.Println("connected to Mongo.")
}

func CloseMongo() {
    MongoSession.Close()
    fmt.Println("disconnected from Mongo.")
}

func WriteMongo(bid_req *openrtb.Request){
    err := MongoC.Insert(bid_req)
    if err != nil {
        log.Fatal(err)
    }
}

func WriteCompressedMongo(bid_req *openrtb.Request){
    buf, err := json.Marshal(bid_req)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return
    }
    fmt.Println(string(buf))
    var b bytes.Buffer
    w := zlib.NewWriter(&b)
    w.Write(buf)
    w.Close()
    fmt.Println(b.Bytes())
/*
    err = MongoC.Insert(bid_req)
    if err != nil {
        log.Fatal(err)
    }
*/
}

func WritePersonMongo(person *Person){
    err := MongoC.Insert(person)
    if err != nil {
        log.Fatal(err)
    }
}
