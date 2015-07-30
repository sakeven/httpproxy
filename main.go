package main

import (
    "httpproxy/cache"
    "httpproxy/proxy"
    "log"
    "net/http"
)

func main() {
    pxy := proxy.NewProxyServer()
    web := proxy.NewWebServer()
    proxy.RegisterCacheBox(cache.NewCacheBox(":6379", ""))
    go http.ListenAndServe(web.Port, web)
    log.Println("begin proxy")
    log.Fatal(pxy.ListenAndServe())
}
