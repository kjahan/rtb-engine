package writer

import (
    "flag"
    "fmt"
    elastigo "github.com/mattbaird/elastigo/lib"
    "log"
    "encoding/json"
    "github.com/bsm/openrtb"
)

var (
	host *string = flag.String("host", "localhost", "Elasticsearch Host")
)

func ConnectES() *elastigo.Conn {
    c := elastigo.NewConn()
    log.SetFlags(log.LstdFlags)
    flag.Parse()
    //set the Elasticsearch Host to Connect to
    c.Domain = *host
    return c
}

func WriteES(c *elastigo.Conn, bid_req *openrtb.Request, id string) {
    //index a document
    response, err := c.Index("rtb", "bid_request", id, nil, bid_req)
    c.Flush()
    if err != nil {
        log.Println("error during serach:" + err.Error())
        log.Fatal(err)
    }
    log.Printf("Index OK: %v", response.Ok)
}

func SearchES(c *elastigo.Conn) {
    // Search Using Raw json String
    searchJson := `{
        "query" : {
	    "term" : { "Name" : "wanda" }
	}
    }`
    out, err := c.Search("testindex", "user", nil, searchJson)
    if err != nil {
	log.Println("error during serach:" + err.Error())
	log.Fatal(err)
    }
    //try to marshalig to MyUser type*
    var u MyUser
    bytes, err := out.Hits.Hits[0].Source.MarshalJSON()
    if err != nil {
	log.Fatalf("err calling marshalJson:%v", err)
    }
    json.Unmarshal(bytes, &u)
    fmt.Println(fmt.Sprintf("%#v", bytes))
    fmt.Println(fmt.Sprintf("%#v", u))
    fmt.Println(fmt.Sprintf("%#v", out.Hits.Hits[0].Source))
}

type MyUser struct {
	Name string
	Age  int
}
