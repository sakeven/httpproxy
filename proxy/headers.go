package proxy

import (
	"net/http"
)

// CopyHeaders copy headers from source to destination.
// Nothing would be returned.
// 拷贝HTTP头
func CopyHeaders(rw, resp http.Header) {
	for key, values := range resp {
		for _, value := range values {
			rw.Add(key, value)
		}
	}
}

// ClearHeaders clear headers.
// 清理HTTP头
func ClearHeaders(headers http.Header) {
	for key, _ := range headers {
		headers.Del(key)
	}
}

// RmProxyHeaders remove Hop-by-hop headers.
// 删除http逐跳头
func RmProxyHeaders(req *http.Request) {
	req.RequestURI = ""
	req.Header.Del("Proxy-Connection")
	req.Header.Del("Connection")
	req.Header.Del("Keep-Alive")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("TE")
	req.Header.Del("Trailers")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Upgrade")
}
