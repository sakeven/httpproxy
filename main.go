package main

import (
    "httpproxy/proxy"
    "log"
    "net/http"
)

func main() {
    pxy := proxy.NewProxyServer()
    web := proxy.NewWebServer()

    go http.ListenAndServe(web.Port, web)
    log.Println("begin proxy")
    log.Fatal(pxy.ListenAndServe())
}
