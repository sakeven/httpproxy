package main

import (
	"log"
	"net/http"

	"github.com/sakeven/httpproxy/proxy"
)

func main() {
	pxy := proxy.NewServer()
	web := proxy.NewWebServer()

	go http.ListenAndServe(web.Port, web)
	log.Println("Begin proxy")
	log.Fatal(pxy.ListenAndServe())
}
