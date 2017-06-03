package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
  "gopkg.in/headzoo/surf.v1"
  "github.com/headzoo/surf/browser"
)

func setCookieHeader(siteCookies []*http.Cookie) string {
  cookies := ""
  for _, cookie := range siteCookies {
    cookies += cookie.String() + ";"
  }

  return cookies
}

func login(cookies []*http.Cookie, token string) {
  var jsonStr = []byte(`
    {
      "associateAccount":"false",
      "email":"alexmnewsdsdman95@gmail.com",
      "pin":"53510560"
    }
  `)

  req, err := http.NewRequest("POST", "https://www.puregym.com/api/members/login/", bytes.NewBuffer(jsonStr))

  req.Header = http.Header{
    "Cookie": []string{setCookieHeader(cookies)},
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
    fmt.Println("error", err)
    return
  }

  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
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

func main() {
  login(getCookies())
}
