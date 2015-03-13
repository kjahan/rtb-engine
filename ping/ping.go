package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I'm still alive!")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":5002", nil)
}
