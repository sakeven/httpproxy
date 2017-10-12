// Package proxy implements a http proxy.
//
// Support GET, POST, CONNECT method and so on.
// Support proxy auth and web management.
// Support web cache.
package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/sakeven/httpproxy/cache/redis"
)

// Server is a server of proxy.
type Server struct {
	// User records user's name
	Tr   *http.Transport
	User string
}

// NewServer returns a new proxyserver.
func NewServer() *http.Server {
	if cnfg.Cache {
		redis.Register(":6379", "")
	}

	return &http.Server{
		Addr:           cnfg.Port,
		Handler:        &Server{Tr: &http.Transport{Proxy: http.ProxyFromEnvironment, DisableKeepAlives: true}},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// ServeHTTP will be automatically called by system.
// ProxyServer implements the Handler interface which need ServeHTTP.
func (proxy *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Debugf("Panic: %v\n", err)
			fmt.Fprintf(rw, fmt.Sprintln(err))
		}
	}()

	// log.Debugf("Host := %v", req.URL.Host)

	if proxy.Auth(rw, req) {
		return
	}

	proxy.ReverseHandler(req)

	if proxy.Ban(rw, req) {
		return
	}

	if req.Method == "CONNECT" {
		proxy.HTTPSHandler(rw, req)
	} else if cnfg.Cache == true && req.Method == "GET" {
		proxy.CacheHandler(rw, req)
	} else {
		proxy.HTTPHandler(rw, req)
	}
}

// HTTPHandler handles http connections.
// 处理普通的http请求
func (proxy *Server) HTTPHandler(rw http.ResponseWriter, req *http.Request) {
	log.Infof("%v is sending request %v %v \n", proxy.User, req.Method, req.URL.Host)
	RmProxyHeaders(req)

	resp, err := proxy.Tr.RoundTrip(req)
	if err != nil {
		log.Errorf("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) //写入响应状态

	nr, err := ioCopy(rw, resp.Body)
	if err != nil && err != io.EOF {
		log.Errorf("%v got an error when copy remote response to client. %v\n", proxy.User, err)
		return
	}
	log.Infof("%v copied %v bytes from %v.\n", proxy.User, nr, req.URL.Host)
}

// HTTP200 http 200 response
var HTTP200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

// HTTPSHandler handles any connection which need connect method.
// 处理https连接，主要用于CONNECT方法
func (proxy *Server) HTTPSHandler(rw http.ResponseWriter, req *http.Request) {
	log.Infof("%v tried to connect to %v", proxy.User, req.URL.Host)

	hj, _ := rw.(http.Hijacker)
	client, _, err := hj.Hijack() //获取客户端与代理服务器的tcp连接
	if err != nil {
		log.Errorf("%v failed to get Tcp connection of %s \n", proxy.User, req.RequestURI)
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}

	remote, err := net.Dial("tcp", req.URL.Host) //建立服务端和代理服务器的tcp连接
	if err != nil {
		log.Errorf("%v failed to connect %v\n", proxy.User, req.RequestURI)
		// TODO write error msg.
		client.Close()
		return
	}

	client.Write(HTTP200)

	go ioCopy(remote, client)
	go ioCopy(client, remote)
}

func ioCopy(dst io.Writer, src io.ReadCloser) (nr int, err error) {
	var buf = GetBuf()
	defer func() {
		PutBuf(buf)
		src.Close()
	}()

	var rerr, werr error
	var n int
	for {
		n, rerr = src.Read(buf)
		if n > 0 {
			_, werr = dst.Write(buf[:n])
			if flusher, ok := dst.(http.Flusher); ok {
				flusher.Flush()
			}
			nr += n
		}

		err = rerr
		if werr != nil {
			err = werr
		}
		if err != nil {
			return
		}
	}
}
