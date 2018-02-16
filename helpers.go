package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// StringifyCookies any given cookies
func StringifyCookies(siteCookies []*http.Cookie) string {
	cookies := ""
	for _, cookie := range siteCookies {
		cookies += cookie.String() + ";"
	}

	return cookies
}

// Getenv Helper to get Env variables whilst providing a default value
func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// FormatRequest generates ascii representation of a request https://medium.com/doing-things-right/pretty-printing-http-requests-in-golang-a918d5aaa000
func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
