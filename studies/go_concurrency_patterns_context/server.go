//article: 
package main

import (
    "net/http"
    "time"
    err "errors"
    userip "service/userip"
)

func main() {
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
    //ctx is the Context for this handler. Callinc cancel closes the
    //ctx.Done channel, which is the cancellation signal for requests
    //started by this handler
    var (
        ctx context.Context
        cancel context.CancelFunc
    )

    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    /*ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers, each
    * with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or
    * "mius"), "ms", "s", "m", "h"
    */
    /*FormValue returns the first value for the named component of the query. POST and PUT body parameters take
    * precedence over URL query string values. FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
    * any erros returned by these functions. If key is not present, FormValue returns the empty string. To access
    * multiple values of the same key, call ParseForm and then inspect Request.Form directly.
    */
    if err == nil {
        //The reqeust has a timeout, so create a context that is
        //canceled automatically when the timeout expires
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        ctx, cancel = context.WithCancel(context.Background())
    }
    defer cancel() //Cancel ctx as soon as handleSearch returns

    //Check the search query
    query := req.FormValue("q")
    if query == "" {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    //Store the user IP in ctx for use by code in other packages
    userIP, err := userip.FromRequest(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ctx = userip.NewContext(ctx, userIP)

    //Run the Google search and print the results.
    start := time.Now()
    results, err := google.Search(ctx, query)
    elapsed := time.Since(start)
}
