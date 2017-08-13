package proxy

import (
	"net/http"
	"strings"
)

func (proxy *ProxyServer) Ban(rw http.ResponseWriter, req *http.Request) bool {
	if len(cnfg.GFWList) > 0 {
		return proxy.ban(rw, req)
	}

	return false
}

func (proxy *ProxyServer) ban(rw http.ResponseWriter, req *http.Request) bool {
	for _, gfwlist := range cnfg.GFWList { //屏蔽列表，检查访问对象是否被屏蔽
		if strings.Index(req.RequestURI, gfwlist) != -1 && gfwlist != "" {
			log.Infof("%s try to visit forbidden website %s", proxy.User, req.URL.Host)
			http.Error(rw, "Forbid", 403)
			return true
		}
	}

	return false
}
