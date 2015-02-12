package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"regexp"
)

var staticRegex = regexp.MustCompile("^/static/(css|js)/[a-zA-Z0-9.-]+$")

func postsHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Path[len("/posts/"):])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	doc, err := GetPage(page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	articles := ExtractArticles(doc)
	resp, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "application/json")
	_, _ = w.Write(resp)
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
	if !staticRegex.MatchString(r.URL.Path) {
		log.Println("Invalid file request: " + r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file := r.URL.Path[len("/static"):]
	log.Println("Serving " + file)
	http.ServeFile(w, r, "public/" + file)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index.html")
	http.ServeFile(w, r, "./index.html")
}

func Serve() {
	http.HandleFunc("/posts/", postsHandler)
	http.HandleFunc("/static/", resourcesHandler)
	http.HandleFunc("/", indexHandler)
	log.Println("Running server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
  Serve()
}
