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
	log.Fatal(pxy.ListenAndServe())
}
