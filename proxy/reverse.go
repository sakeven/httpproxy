package proxy

import (
    "net/http"
)

//ReverseHandler handles request for reverse proxy.
//处理反向代理请求
func (proxy *ProxyServer) ReverseHandler(req *http.Request) {
    if cnfg.Reverse == true { //用于反向代理
        proxy.ReverseHandler(req)
    }
}

//ReverseHandler handles request for reverse proxy.
//处理反向代理请求
func (proxy *ProxyServer) reverseHandler(req *http.Request) {
    req.Host = cnfg.Proxy_pass
    req.URL.Host = req.Host
    req.URL.Scheme = "http"
    log.Debug("%v", req.RequestURI)
}
