package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "bytes"
)

func login() {
  var jsonStr = []byte(`{"associateAccount":"false","email":"alexmnewman95@gmail.com","pin":"5351sd0560"}`)
  req, err := http.NewRequest("POST", "https://www.puregym.com/api/members/login/", bytes.NewBuffer(jsonStr))

  req.Header.Add("Cookie", " __cfduid=d300aa5b5c20903fcc47688673f370f201493824932; CookieNotification=; raygun4js-userid=2573633e-abb4-f56e-1246-7918e41d1af4; ARRAffinity=3dd5653018c8bf8bc25af4fce62e6eb7cd7a2075a1cae7e31ac297d81db86e32; ex-sess=495cf1e4-0c70-4d05-b6b3-c8804192fe1b; __RequestVerificationToken=RraCbmNwkW8Uai5gcxfVkcfPf5Ov-Gs8p8k059y8LdBSnYPYaND8y7SC_B0A3s-IdvdLVVu_3L-EJDsLrW_eFhsa3mY1; _gat_UA-9256723-1=1; _vwo_uuid_v2=AE6347DF3C64A6CDD4CFE8F941135B4F|dd5a1dbb15b097ec89b62209ee4461dc; _ga=GA1.2.388286209.1493824935; _gid=GA1.2.1293107821.1496404552; _tq_id.TV-361818-1.e6fc=68e0c24901f48154.1493824935.0.1496404749..")

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
  req.Header.Add("__RequestVerificationToken", "wsabMJop7cPXN3vFq2RdosGpJlek8tuFLHbCIMlirR2UJT1JEpm09ExCqiL4wLRIQSWKMWm5xWMjVYjlF81fB4BI0SY1")
  req.Header.Set("DNT", "1")




  client := &http.Client{}
  resp, err := client.Do(req)



  if nil!= err {
    fmt.Println("error", err)
    return
  }

  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("fuck", req)
  fmt.Println("response Body:", string(body))
}

func main() {
  fmt.Printf("Hello, world.\n")
  login()
}
