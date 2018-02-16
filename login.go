package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

var user = Getenv("PURE_USER", "nononosir@gmail.com")
var pin = Getenv("PURE_PIN", "heirsek")

func getHeaders(cookies string, token string) http.Header {
	return http.Header{
		"Cookie":                     []string{cookies},
		"Origin":                     []string{"https://puregym.com"},
		"Host":                       []string{"https://puregym.com"},
		"Accept-Encoding":            []string{"gzip, deflate, br"},
		"Accept-Language":            []string{"en-GB,en-US;q=0.8,en;q=0.6"},
		"X-Requested-With":           []string{"XMLHttpRequest"},
		"Pragma":                     []string{"no-cache"},
		"User-Agent":                 []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3118.0 Safari/537.36"},
		"Upgrade-Insecure-Requests":  []string{"1"},
		"Content-Type":               []string{"text/html; charset=utf-8"},
		"Accept":                     []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
		"Cache-Control":              []string{"no-cache"},
		"Referer":                    []string{"https://www.puregym.com/Login/?ReturnUrl=%2Fmembers%2F"},
		"__RequestVerificationToken": []string{token},
		"DNT": []string{"1"},
	}
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func readMembers(body string) string {
	buf := bytes.NewBufferString(string(body))
	doc, _ := goquery.NewDocumentFromReader(buf)
	el := doc.Find(".heading.heading--level3.secondary-color.margin-none").Text()

	return el
}

func writeData(body string) {
	jsonBlob, _ := ioutil.ReadFile("./output.json")

	type Count struct {
		Date   string
		People string
	}

	group := Count{
		Date:   getTime(),
		People: readMembers(body),
	}
	var counts []Count
	err := json.Unmarshal(jsonBlob, &counts)
	if err != nil {
		fmt.Println("error:", err)
	}

	var appendData = append(counts, group)
	b, _ := json.Marshal(appendData)

	os.Stdout.Write(b)

	ioutil.WriteFile("output.json", b, 0644)
}

func formatData(body string) {
	writeData(body)
}

func getMembers(cookies string, token string) {
	req, err := http.NewRequest("GET", "https://www.puregym.com/members/", nil)
	req.Header = getHeaders(cookies, token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if nil != err {
		fmt.Println("error", err)
		return
	}

	r, err := gzip.NewReader(resp.Body)
	r.Close()
	body, _ := ioutil.ReadAll(r)
	formatData(string(body))
}

func authenticate(siteCookies []*http.Cookie, token string) {
	cookies := StringifyCookies(siteCookies)
	var details = []byte(`
    {
      "associateAccount":"false",
      "email":"` + user + `",
      "pin":"` + pin + `"
    }
  `)

	req, err := http.NewRequest("POST", "https://www.puregym.com/api/members/login/", bytes.NewBuffer(details))
	req.Header = getHeaders(cookies, token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/javascript")

	client := &http.Client{}
	resp, err := client.Do(req)

	if nil != err {
		fmt.Println("Login has failed!", err)
		return
	}

	cookies += StringifyCookies(resp.Cookies())

	getMembers(cookies, token)
}

func getSite() *browser.Browser {
	browser := surf.NewBrowser()
	err := browser.Open("http://puregym.com")
	if err != nil {
		fmt.Println("error", err)
	}

	return browser
}

func getCookies() ([]*http.Cookie, string) {
	browser := getSite()

	token, _ := browser.Dom().Find("[name='__RequestVerificationToken']").Attr("value")
	cookies := browser.SiteCookies()

	return cookies, token
}

// Creates cookies and logs into puregym
func Login() {
	authenticate(getCookies())
}
