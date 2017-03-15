package main

import (
	"log"
	"net/http"

	"httpproxy/proxy"
)

func main() {
	pxy := proxy.NewProxyServer()
	web := proxy.NewWebServer()

	go http.ListenAndServe(web.Port, web)
	log.Println("begin proxy")
	log.Fatal(pxy.ListenAndServe())
}
