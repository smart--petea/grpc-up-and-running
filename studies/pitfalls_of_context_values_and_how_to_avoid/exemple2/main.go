package main

import (
    "net/http"
    "fmt"
    "context"
)

var (
    requestID = 0
)

const (
    port = ":3000"
)

func nextRequestID() int {
    requestID++
    return requestID
}

/*
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

ServeMux implement such a function
*/

func addRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        ctx := context.WithValue(r.Context(), "request_id", nextRequestID())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func requireUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        /*
        user := lookupUser(r)
        if user == nil {
            //No user so redirect to login
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
        */
        user := 1
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func printHi(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hi! Your request ID is:", r.Context().Value("request_id"))
}

func printBye(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Bye! Your request ID is:", r.Context().Value("request_id"))
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "HOME")
}

func main() {
    dashboard := http.NewServeMux()
    dashboard.HandleFunc("/dashboard/hi", printHi)
    dashboard.HandleFunc("/dashboard/bye", printBye)

    mux := http.NewServeMux()
    //ALL routes that start with /dashboard/ require that a
    //user is authenticated using the requireUser middleware
    mux.Handle("/dashboard/", requireUser(dashboard))
    mux.HandleFunc("/", home)

    http.ListenAndServe(port, addRequestID(mux))
}
