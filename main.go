package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var port = Getenv("PORT", "8080")
var test string

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "src/index.html")
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/index.html"))
	type TodoPageData struct {
		Count string
	}
	data := TodoPageData{
		Count: Login(),
	}
	tmpl.Execute(w, data)
	fmt.Println("server started")
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r)
	Login()
}

func main() {
	http.HandleFunc("/members", membersHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
