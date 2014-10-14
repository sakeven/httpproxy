package proxy

import (
	"httpproxy/cache"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var Caches cache.CacheSet

//IsCache checks whether response can be stored as cache
//通过对http头的检测确定是否缓存
func IsCache(resp *http.Response, URI string) bool {
	log.Debug("%v", URI)
	Cache_Control := resp.Header.Get("Cache-Control")
	Content_type := resp.Header.Get("Content-Type")
	log.Debug("Cache-Control: %v", Cache_Control)
	if strings.Index(Cache_Control, "private") != -1 ||
		strings.Index(Cache_Control, "no-store") != -1 ||
		strings.Index(Content_type, "application") != -1 ||
		strings.Index(Content_type, "video") != -1 ||
		strings.Index(Content_type, "audio") != -1 ||
		(strings.Index(Cache_Control, "max-age") == -1 &&
			strings.Index(Cache_Control, "s-maxage") == -1 &&
			resp.Header.Get("Etag") == "" &&
			resp.Header.Get("Last-Modified") == "" &&
			(resp.Header.Get("Expires") == "" || resp.Header.Get("Expires") == "0")) {
		log.Debug("False")
		return false
	}
	return true
}

//CacheHandler handles "Get" request
// 只处理get请求的缓存
func (proxy *ProxyServer) CacheHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		nr       int
		err      error
		Hit      bool //did hit the cache
		RepCache bool //should be replace old cache
	)
	URI := req.RequestURI
	Cache := Caches[URI]

	if Cache != nil {
		if Cache.Mustverified == false && Cache.Vlidity.After(time.Now().UTC()) {
			Hit = true
		} else {
			if Cache.Verify() {
				Hit = true
			}
		}
		if Hit == false {
			RepCache = true
		}
	}

	remoteinfo := "" // 记录响应是来自本地缓存还是远端服务器

	if Hit {
		log.Debug("Hit %v", URI)
		remoteinfo = "local cache"
		CopyHeaders(rw.Header(), Cache.Header)
		rw.WriteHeader(Cache.StatusCode)
		nr, err = rw.Write([]byte(Cache.Body))
	} else {
		RmProxyHeaders(req)
		resp, err := proxy.Tr.RoundTrip(req)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			log.Debug("it's %v", err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("%v", err)
		}
		if IsCache(resp, URI) {
			if RepCache == false {
				Caches[URI] = &cache.Cache{URI: URI}
				Cache = Caches[URI]
			}
			log.Debug("Store Cache %v", URI)
			Cache.CopyHeaders(resp.Header)
			err = Cache.SetCache(resp.StatusCode, string(b))
			if err != nil {
				log.Error("%v", err)
			}
		} else {
			Caches.Delete(URI)
		}

		remoteinfo = "remote"
		CopyHeaders(rw.Header(), resp.Header)
		rw.WriteHeader(resp.StatusCode)
		nr, err = rw.Write([]byte(b))
	}

	if err != nil && err != io.EOF {
		log.Error("%s got an error when copy %s response to client.%v\n", proxy.User, remoteinfo, err)
		return
	}
	log.Info("%s Copied %d bytes from %s %s.\n", proxy.User, nr, remoteinfo, req.URL.Host)
	return
}

//ExistCache check wether specific uri cache exists
//检查缓存是否已存在
func ExistCache(uri string) bool {
	if Caches[uri] == nil {
		return false
	} else {
		return true
	}
}

//CheckCaches evey certian minutes check whether cache is out of date, if yes release it.
//每隔特定时间释放过期缓存
func CheckCaches() {
	for {
		time.Sleep(time.Duration(cnfg.CacheTimeout) * time.Minute)
		for key, Cache := range Caches {
			if Cache != nil && Cache.Verify() == false {
				Caches.Delete(key)
			}
		}
	}
}
