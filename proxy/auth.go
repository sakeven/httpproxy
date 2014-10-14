package proxy

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var HTTP_407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Secure Proxys\"\r\n\r\n")

//Auth provides basic authorizaton for proxy server.
//验证用户登入，如果不能登入将返回一个error
func (proxy *ProxyServer) Auth(rw http.ResponseWriter, req *http.Request) (string, error) {

	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)

	if auth == "" {
		NeedAuth(rw, HTTP_407)
		return "", errors.New("Need Proxy Authorization!")
	}
	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		log.Debug("when decoding %v, got an error of %v", auth, err)
		return "", errors.New("Fail to decoding Proxy-Authorization")
	}

	var user, passwd string

	userPasswdPair := strings.Split(string(data), ":")
	if len(userPasswdPair) != 2 {
		NeedAuth(rw, HTTP_407)
		return "", errors.New("Fail to log in")
	} else {
		user = userPasswdPair[0]
		passwd = userPasswdPair[1]
	}
	if Check(user, passwd) == false {
		NeedAuth(rw, HTTP_407)
		return "", errors.New("Fail to log in")
	}
	return user, nil
}

// 代理服务器需要验证，发送验证响应
func NeedAuth(rw http.ResponseWriter, challenge []byte) error {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		return errors.New("Fail to get Tcp connection of Client")
	}
	defer Client.Close()

	Client.Write(challenge)
	return nil
}

// 验证用户名和密码
func Check(user, passwd string) bool {
	if user != "" && passwd != "" && cnfg.User[user] == passwd {
		return true
	} else {
		return false
	}
}
