// Package proxy implements a http proxy.
//
// Support GET,POST,CONNECT method and so on.
// Support proxy auth and web management.
// Support web cache.
package proxy

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type ProxyServer struct {
	// User records user's name
	Tr   *http.Transport
	User string
}

// NewProxyServer returns a new proxyserver.
func NewProxyServer() *http.Server {
	return &http.Server{
		Addr:           cnfg.Port,
		Handler:        &ProxyServer{Tr: &http.Transport{Proxy: http.ProxyFromEnvironment}},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

//ServeHTTP will be automatically called by system.
//ProxyServer implements the Handler interface which need ServeHTTP.
func (proxy *ProxyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var err error

	log.Debug("Host := %v", req.URL.Host)

	if cnfg.Reverse == false && cnfg.Auth == true { //代理服务器登入认证
		if proxy.User, err = proxy.Auth(rw, req); err != nil {
			log.Error("%v", err)
			return
		}
	} else {
		proxy.User = "Anonymous"
	}

	if cnfg.Reverse == true { //用于反向代理
		proxy.ReverseHandler(req)
	}

	for _, gfwlist := range cnfg.GFWList { //屏蔽列表，检查访问对象是否被屏蔽
		if strings.Index(req.RequestURI, gfwlist) != -1 && gfwlist != "" {
			log.Info("%s try to visit forbidden website %s", proxy.User, req.URL.Host)
			http.Error(rw, "Forbid", 403)
			return
		}
	}
	if req.Method == "CONNECT" {
		proxy.HttpsHandler(rw, req)
	} else if cnfg.Cache == true && req.Method == "GET" {
		proxy.CacheHandler(rw, req)
	} else {
		proxy.HttpHandler(rw, req)
	}
}

//ReverseHandler handles request for reverse proxy.
//处理反向代理请求
func (proxy *ProxyServer) ReverseHandler(req *http.Request) {
	req.Host = cnfg.Proxy_pass
	req.URL.Host = req.Host
	req.URL.Scheme = "http"
	log.Debug("%v", req.RequestURI)
}

//HttpHandler handles http connections.
//处理普通的http请求
func (proxy *ProxyServer) HttpHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info("%v is sending request %v %v \n", proxy.User, req.Method, req.URL.Host)
	RmProxyHeaders(req)

	resp, err := proxy.Tr.RoundTrip(req)
	if err != nil {
		log.Error("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) //写入响应状态

	nr, err := io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		log.Error("%v got an error when copy remote response to client.%v\n", proxy.User, err)
		return
	}
	log.Info("%v Copied %v bytes from %v.\n", proxy.User, nr, req.URL.Host)
}

var HTTP_200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

// HttpsHandler handles any connection which need connect method.
// 处理https连接，主要用于CONNECT方法
func (proxy *ProxyServer) HttpsHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info("%v tried to connect to %v", proxy.User, req.URL.Host)

	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack() //获取客户端与代理服务器的tcp连接
	if err != nil {
		log.Error("%v failed to get Tcp connection of \n", proxy.User, req.RequestURI)
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}

	Remote, err := net.Dial("tcp", req.URL.Host) //建立服务端和代理服务器的tcp连接
	if err != nil {
		log.Error("%v failed to connect %v\n", proxy.User, req.RequestURI)
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}

	Client.Write(HTTP_200)

	go copyRemoteToClient(proxy.User, Remote, Client)
	go copyRemoteToClient(proxy.User, Client, Remote)
}

func copyRemoteToClient(User string, Remote, Client net.Conn) {
	defer Client.Close()
	nr, err := io.Copy(Remote, Client)
	if err != nil && err != io.EOF {
		log.Error("%v got an error when handles CONNECT %v\n", User, err)
		return
	}
	log.Info("%v transport %v bytes betwwen %v and %v.\n", User, nr, Remote.RemoteAddr(), Client.RemoteAddr())
}
