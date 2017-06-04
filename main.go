package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
  "gopkg.in/headzoo/surf.v1"
  "github.com/headzoo/surf/browser"
  "compress/gzip"
  "github.com/PuerkitoBio/goquery"
  "time"
  "encoding/json"
  "os"
)

func stringifyCookies(siteCookies []*http.Cookie) string {
  cookies := ""
  for _, cookie := range siteCookies {
    cookies += cookie.String() + ";"
  }

  return cookies
}

func login(siteCookies []*http.Cookie, token string) {
  cookies := stringifyCookies(siteCookies)
  var jsonStr = []byte(`
    {
      "associateAccount":"false",
      "email":"alexmnewman95@gmail.com",
      "pin":"53510560"
    }
  `)

  req, err := http.NewRequest("POST", "https://www.puregym.com/api/members/login/", bytes.NewBuffer(jsonStr))

  req.Header = http.Header{
    "Cookie": []string{cookies},
    "Origin": []string{"https://puregym.com"},
    "Accept-Encoding": []string{"gzip, deflate, br"},
    "Accept-Language": []string{"en-GB,en-US;q=0.8,en;q=0.6"},
    "X-Requested-With": []string{"XMLHttpRequest"},
    "Pragma": []string{"no-cache"},
    "User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3118.0 Safari/537.36"},
    "Content-Type": []string{"application/json"},
    "Accept": []string{"application/json, text/javascript"},
    "Cache-Control": []string{"no-cache"},
    "Referer": []string{"https://www.puregym.com/Login/?ReturnUrl=%2Fmembers%2F"},
    "__RequestVerificationToken": []string{token},
    "DNT": []string{"1"},
  }

  client := &http.Client{}
  resp, err := client.Do(req)

  if nil!= err {
    fmt.Println("Login has failed!", err)
    return
  }

  cookies += stringifyCookies(resp.Cookies())

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

func getMembers(cookies string, token string) {
  req, err := http.NewRequest("GET", "https://www.puregym.com/members/", nil)

  req.Header = http.Header{
    "Cookie": []string{cookies},
    "Origin": []string{"https://puregym.com"},
    "Host": []string{"https://puregym.com"},
    "Accept-Encoding": []string{"gzip, deflate, br"},
    "Accept-Language": []string{"en-GB,en-US;q=0.8,en;q=0.6"},
    "X-Requested-With": []string{"XMLHttpRequest"},
    "Pragma": []string{"no-cache"},
    "User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3118.0 Safari/537.36"},
    "Upgrade-Insecure-Requests": []string{"1"},
    "Content-Type": []string{"text/html; charset=utf-8"},
    "Accept": []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
    "Cache-Control": []string{"no-cache"},
    "Referer": []string{"https://www.puregym.com/Login/"},
    "__RequestVerificationToken": []string{token},
    "DNT": []string{"1"},
  }

  client := &http.Client{}
  resp, err := client.Do(req)

  if nil!= err {
    fmt.Println("error", err)
    return
  }

  r, err := gzip.NewReader(resp.Body)
  r.Close()
  body, _ := ioutil.ReadAll(r)
  formatData(string(body))
}

func getTime() string {
  return time.Now().Format(time.RFC3339)
}

func readMembers(body string) string {
  buf := bytes.NewBufferString(string(body))
  doc, _ := goquery.NewDocumentFromReader(buf)
  el := doc.Find(".heading.heading--level3.secondary-color.margin-none").Text()

  return el;
}

func formatData(body string) {
  var m = make(map[string]string)
  m[getTime()] = readMembers(body)
  fmt.Println(m)
  data, _ := json.Marshal(m)
  ioutil.WriteFile("output.json", data, 0644)
  os.Stdout.Write(data)
}

func main() {
  login(getCookies())
}
