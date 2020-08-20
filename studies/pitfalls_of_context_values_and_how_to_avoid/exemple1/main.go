package main

import (
    "net/http"
    "log"
    "context"
    "fmt"
)

var (
    requestID = 0
)

const (
    port = ":9998"
)

func nextRequestID() int {
    requestID = requestID + 1
    log.Println(requestID)
    return requestID
}

func addRequestID(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("a")
        ctx := context.WithValue(r.Context(), "request_id", nextRequestID())
        log.Println("b")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello request id = %d", r.Context().Value("request_id"))
}

func main() {
    http.HandleFunc("/", addRequestID(handler))
    log.Fatal(http.ListenAndServe(port, nil))
}
