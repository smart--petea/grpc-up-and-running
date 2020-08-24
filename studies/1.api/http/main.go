package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

func main() {
    Articles = []Article {
        Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2", Desc: "Article Description 2", Content: "Article Content 2"},
    }

    handleRequests()
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    /*Vars returns the route variables for the current request, if any*/
    vars := mux.Vars(r)
    key := vars["id"]

    for _, article  := range Articles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
            return
        }
    }

    http.NotFound(w, r)
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

type Article struct {
    Id string `json: Id"`
    Title string `json:"Title"`
    Desc string `json:"Desc"`
    Content string `json:"Content"`
}

var Articles []Article

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

/*
func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/articles", returnAllArticles)
    log.Fatal(http.ListenAndServe(":10000", nil))
}
*/

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/all", returnAllArticles)
    myRouter.HandleFunc("/article/{id}", returnSingleArticle)

    //finally, instead of passing in nil, we want
    //to pass in our newly created router as the second argument
    log.Fatal(http.ListenAndServe(":10001", myRouter))
}
