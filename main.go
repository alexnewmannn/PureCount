package main

import (
  "fmt"
  "net/http"
)

var port = Getenv("PORT", "8080")

func serveFile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  http.ServeFile(w, r, "src/index.html")
}

func handler(w http.ResponseWriter, r *http.Request) {
  serveFile(w, r)
  fmt.Println("server started")
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
  serveFile(w, r)
  Login()
}

func main() {
  http.HandleFunc("/members", membersHandler)
  http.HandleFunc("/", handler)
  http.ListenAndServe("localhost:"+port, nil)
}
