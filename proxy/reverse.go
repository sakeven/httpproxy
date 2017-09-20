package proxy

import (
	"net/http"
)

// ReverseHandler handles request for reverse proxy.
func (proxy *Server) ReverseHandler(req *http.Request) {
	if cnfg.Reverse == true { //用于反向代理
		proxy.reverseHandler(req)
	}
}

// reverseHandler handles request for reverse proxy.
func (proxy *Server) reverseHandler(req *http.Request) {
	req.Host = cnfg.ProxyPass
	req.URL.Host = req.Host
	req.URL.Scheme = "http"
	log.Debugf("%v", req.RequestURI)
}
