package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
  "gopkg.in/headzoo/surf.v1"
  "github.com/PuerkitoBio/goquery"
)

func createHeaders(test []*http.Cookie) string {
  var lol = ""
  for _, item := range test {
    lol += item.String() + ";"
  }

  return lol
}

func login(test []*http.Cookie) {
  var jsonStr = []byte(`{"associateAccount":"false","email":"alexmnewdman95@gmail.com","pin":"32250194"}`)
  req, err := http.NewRequest("POST", "https://www.puregym.com/api/members/login/", bytes.NewBuffer(jsonStr))

  createHeaders(test)

  // var lol = ""
  // lol += test

  // fmt.Println(createHeaders(test))

  req.Header.Add("Cookie", createHeaders(test))

  req.Header.Set("Origin", "https://www.puregym.com")
  req.Header.Set("Accept-Encoding", "gzip, deflate, br")
  req.Header.Set("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")
  req.Header.Set("X-Requested-With", "XMLHttpRequest")
  req.Header.Set("Connection", "keep-alive")
  req.Header.Set("Pragma", "no-cache")
  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3118.0 Safari/537.36")
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept", "application/json, text/javascript")
  req.Header.Set("Cache-Control", "no-cache")
  req.Header.Set("Referer", "https://www.puregym.com/Login/?ReturnUrl=%2Fmembers%2F")
  // This is retrieved from the DOM
  req.Header.Add("__RequestVerificationToken", "Ht3u55kSYdlooL3H_LPWIcX3Fl51bHKPr8y97w6L5WmjOp74IbRmA2LtmgVvJ0IbPaIwkEZkXjy1nWIeVpuaq9aMCfA1")
  req.Header.Set("DNT", "1")

  client := &http.Client{}
  resp, err := client.Do(req)

  if nil!= err {
    fmt.Println("error", err)
    return
  }

  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
}

func getCookies() []*http.Cookie {
  browser := surf.NewBrowser()
  err := browser.Open("http://puregym.com")
  browser.Dom()

  // browser.Dom().Find("input").Each(func(_ int, s *goquery.Selection) {
  //   fmt.Println(s.OuterHtml())
  // })
  // fmt.Println(test)
  cookies := browser.SiteCookies()
  if err != nil {
    fmt.Println("error", err)
  }
  return cookies
}


func main() {
  fmt.Printf("Hello, world.\n")
  // var test = getCookies()

  login(getCookies())
}
