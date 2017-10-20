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
  "os"
  "encoding/json"
  "github.com/garyburd/redigo/redis"
  "reflect"
  "strconv"
)

var redisDb redis.Conn
var port = os.Getenv("PORT")
var user = os.Getenv("PURE_USER")
var pin = os.Getenv("PURE_PIN") // handle not having these

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
      "email":"` + user + `",
      "pin":"` + pin + `"
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
  if err != nil {
    fmt.Println("error:", err)
  } else {
    r.Close()
    body, _ := ioutil.ReadAll(r)
    formatData(string(body))
  }
}

func getTime() int64 {
  return time.Now().Unix()
}

func readMembers(body string) string {
  buf := bytes.NewBufferString(string(body))
  doc, _ := goquery.NewDocumentFromReader(buf)
  el := doc.Find(".heading.heading--level3.secondary-color.margin-none").Text()

  return el;
}

func formatData(body string) {
  writeData(body)
}

func writeData(body string) {
  jsonBlob, _ := ioutil.ReadFile("./output.json")

  type Animal struct {
    Date  string
    People string
  }

  // v := int64(getTime())
  time := strconv.FormatInt(getTime(), 10)

  group := Animal{
    Date:     time,
    People:   readMembers(body),
  }
  var animals []Animal
  err := json.Unmarshal(jsonBlob, &animals)
  if err != nil {
    fmt.Println("error unmarshal:", err)
  }

  var test = append(animals, group)
  b, _ := json.Marshal(test)

  os.Stdout.Write(b)
  n, redisErr := redisDb.Do("ZADD", "members", getTime(), time + ":" +  readMembers(body))
  fmt.Println(n)
  if redisErr != nil {
    fmt.Println("error zadd:", redisErr)
  }

  scanRedis()
  ioutil.WriteFile("output.json", b, 0644)
}

func scanRedis() {
  // values, err := redisDb.Do("ZRANGE", "members", 0, -1)
  members, err := redis.Strings(redisDb.Do("ZRANGE", "members", 0, -1))

  if err != nil {
    fmt.Println("error zrange:", err)
  }
  fmt.Println(members)
}

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  http.ServeFile(w, r, "output.json")
}

func connectRedis() {
  conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
  if err != nil {
    fmt.Println("error:", err)
  }
  defer conn.Close()
}

func worker() {
  for {
    login(getCookies())
    time.Sleep(time.Second * 60)
  }
}

func main() {
  go worker()
  // connectRedis()
  var err error
  redisDb, err = redis.DialURL(os.Getenv("REDIS_URL"))
  // conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
  if err != nil {
    fmt.Println("error:", err)
  }
  fmt.Println(reflect.TypeOf(redisDb).Kind())
  http.HandleFunc("/", handler)
  http.ListenAndServe("localhost:"+port, nil)
}
