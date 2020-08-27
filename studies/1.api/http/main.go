package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "io/ioutil"
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
    log.Printf("returnSingleArticle %s %s\n", r.URL.Path, r.Method)
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

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    log.Printf("deleteArticle %s\n", r.URL.Path)
    vars := mux.Vars(r)
    id := vars["id"]

    for index, article := range Articles {
        if article.Id == id {
            Articles = append(Articles[:index], Articles[index+1:]...)
            return;
        }
    }

    http.NotFound(w, r)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    log.Printf("createNewArticle %s\n", r.URL.Path)
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var article Article
    err = json.Unmarshal(reqBody, &article)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/all", returnAllArticles)
    myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")

    //finally, instead of passing in nil, we want
    //to pass in our newly created router as the second argument
    log.Fatal(http.ListenAndServe(":10001", myRouter))
}
