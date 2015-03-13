default: deps

deps:
	go get github.com/bsm/openrtb
	go get github.com/lib/pq
	go get github.com/mattbaird/elastigo
	go get gopkg.in/mgo.v2
